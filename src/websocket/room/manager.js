const { customAlphabet } = require('nanoid')
const Room = require('./room')

const nanoid = customAlphabet('1234567890abcdef', 10)

class RoomManager {
    constructor() {
        this.rooms = new Map()
    }

    create(type) {
        const id = nanoid() // todo: make this dns compatible
        const room = new Room(id, type)

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
