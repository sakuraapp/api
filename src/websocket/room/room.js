const { EventEmitter } = require('events')

class Room extends EventEmitter {
    constructor(id) {
        super()

        this.id = id
        this.clients = []
        this.owner = null
    }

    add(client) {
        if (!this.owner) {
            this.owner = client
        }

        if (client.room) {
            client.room.remove(client)
        }

        this.send({
            action: 'add user',
            data: client.profile,
        })

        this.clients.push(client)

        client.room = this
        client.send({
            action: 'join room',
            data: this.getRoomInfo(),
        })

        this.sendRoomState(client)
    }

    remove(client) {
        const i = this.clients.indexOf(client)

        if (i > -1) {
            this.splice(i, 1)
            this.send({
                action: 'remove user',
                data: client.id,
            })

            client.room = null
        }
    }

    getRoomInfo() {
        return {
            owner: this.owner.id,
            users: this.clients.map((client) => client.profile),
        }
    }

    sendRoomState(target) {
        if (!target) {
            target = this
        }

        //target.send()
    }
}

module.exports = Room
