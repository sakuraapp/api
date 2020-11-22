const { EventEmitter } = require('events')
const Opcodes = require('@common/opcodes.json')
const Permissions = require('@common/permissions.json')
const Queue = require('./queue')

class Room extends EventEmitter {
    constructor(id) {
        super()

        this.id = id
        this.ownerId = null
        this.ownerUsername = null
        this.private = false

        this.clients = []
        this.messages = []
        this.invites = []

        this.state = null
        this.queue = new Queue(this)
    }

    get owner() {
        return this.clients.find((client) => client.id === this.ownerId)
    }

    send(data, ignored = []) {
        if (ignored && !Array.isArray(ignored)) ignored = [ignored]

        for (const client of this.clients) {
            if (!ignored.includes(client)) {
                client.send(data)
            }
        }
    }

    add(client) {
        if (!this.ownerId) {
            this.ownerId = client.id
        }

        if (client.id === this.ownerId) {
            this.ownerUsername = client.username
        }

        if (client.room) {
            client.room.remove(client)
        }

        if (this.can(client, 'join room')) {
            const index = this.invites.indexOf(client)

            if (index > -1) {
                this.invites.splice(index, 1)
            }

            this.send({
                op: Opcodes.ADD_USER,
                d: client.profile,
            })

            this.clients.push(client)

            client.room = this
            client.send({
                op: Opcodes.JOIN_ROOM,
                d: {
                    status: 200,
                    room: this.getInfo(),
                    permissions: this.getPermissions(client),
                },
            })

            this.sendPlayerState(client)
        } else {
            client.send({
                op: Opcodes.JOIN_ROOM,
                d: { status: 401 },
            })

            const owner = this.owner

            if (owner) {
                this.owner.send({
                    op: Opcodes.ROOM_JOIN_REQUEST,
                    d: client.profile,
                })
            }
        }
    }

    remove(client) {
        const i = this.clients.indexOf(client)

        if (i > -1) {
            this.clients.splice(i, 1)
            this.send({
                op: Opcodes.REMOVE_USER,
                d: client.id,
            })

            client.room = null
        }

        client.send({
            op: Opcodes.LEAVE_ROOM,
            d: this.id,
        })
    }

    getInfo() {
        return {
            id: this.id,
            owner: this.owner.id,
            users: this.clients.map((client) => client.profile),
            queue: this.queue.items,
        }
    }

    getPublicInfo() {
        const currentItem = this.private ? null : this.queue.currentItem

        return {
            id: this.id,
            owner: this.ownerUsername,
            private: this.private,
            currentItem,
        }
    }

    getState() {
        const state = Object.assign({}, this.state)

        if (state.playing) {
            state.currentTime = this.getCurrentTime()
        }

        return state
    }

    getCurrentTime() {
        console.log(
            this.state.currentTime,
            (new Date().getTime() - this.state.playbackStart) / 1000
        )
        return (
            this.state.currentTime +
            (new Date().getTime() - this.state.playbackStart) / 1000
        )
    }

    getPermissions(client) {
        if (client.id === this.ownerId) {
            return Object.values(Permissions)
        }

        return [Permissions.QUEUE_ADD]
    }

    setPlaying(val, send = true) {
        if (!val) {
            if (this.state.playing) {
                this.state.currentTime = this.getCurrentTime()
            }
        } else {
            this.state.playbackStart = new Date().getTime()
        }

        this.state.playing = val

        if (send) this.sendPlayerState()
    }

    sendPlayerState(target, ignored) {
        if (!target) {
            target = this
        }

        target.send(
            {
                op: Opcodes.PLAYER_STATE,
                d: this.getState(),
            },
            ignored
        )
    }

    sendMessage(content, client) {
        const message = {
            content,
            author: client.id,
            time: new Date().getTime(),
        }

        message.id = this.messages.push(message)

        this.send(
            {
                op: Opcodes.SEND_MESSAGE,
                d: message,
            },
            client
        )

        return message.id
    }

    getClientById(id) {
        return this.clients.find((client) => client.id === id)
    }

    getClientBySocketId(socketId) {
        return this.clients.find((client) => client.socketId === socketId)
    }

    hasUser(id) {
        return this.getClientById(id) !== null
    }

    hasPermission(name, client) {
        return this.can(client, name)
    }

    can(client, action) {
        switch (action) {
            case 'join room':
                return (
                    !this.private ||
                    client.id === this.ownerId ||
                    this.invites.includes(client)
                )
                break
            default:
                const perms = this.getPermissions(client)

                return perms.includes(Permissions[action])
        }
    }

    findClientsWithPermissions(actions) {
        const clients = []
        const permissions = actions.map((action) => Permissions[action])

        for (const client of this.clients) {
            const perms = this.getPermissions(client)
            let valid = true

            for (const permission of permissions) {
                if (!perms.includes(permission)) {
                    valid = false

                    break
                }
            }

            if (valid) {
                clients.push(client)
            }
        }

        return clients
    }
}

module.exports = Room
