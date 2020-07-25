const shortid = require('shortid')
const Room = require('./room')

class RoomManager {
    constructor() {
        this.rooms = new Map()
    }

    create() {
        const id = shortid.generate()
        const room = new Room(id)

        room.on('destroy', () => {
            this.rooms.delete(id)
        })

        this.rooms.set(id, room)

        return room
    }

    find(id) {
        return this.rooms.get(id)
    }

    findByOwner(id) {
        for (const room of this.rooms.values()) {
            if (room.ownerId === id) {
                return room
            }
        }
    }

    join(id, client) {
        const room = this.find(id)

        if (!room) {
            throw new Error(`Room of ID ${id} doesn't exist.`)
        }

        room.add(client)
    }
}

module.exports = new RoomManager()
