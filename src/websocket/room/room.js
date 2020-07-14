const { EventEmitter } = require('events')
const Opcodes = require('@common/opcodes.json')

class Room extends EventEmitter {
    constructor(id) {
        super()

        this.id = id
        this.clients = []
        this.owner = null
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
        if (!this.owner) {
            this.owner = client
        }

        if (client.room) {
            client.room.remove(client)
        }

        this.send({
            op: Opcodes.ADD_USER,
            d: client.profile,
        })

        this.clients.push(client)

        client.room = this
        client.send({
            op: Opcodes.JOIN_ROOM,
            d: this.getRoomInfo(),
        })

        this.sendPlayerState(client)
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
    }

    getRoomInfo() {
        return {
            id: this.id,
            owner: this.owner.id,
            users: this.clients.map((client) => client.profile),
        }
    }

    sendPlayerState(target) {
        if (!target) {
            target = this
        }

        //target.send()
    }
}

module.exports = Room
