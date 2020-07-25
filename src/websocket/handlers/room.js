const Handler = require('./handler')
const Opcodes = require('@common/opcodes.json')
const roomManager = require('../room/manager')

const handler = new Handler()

handler.on(Opcodes.DISCONNECT, (data, client) => {
    if (client.room) {
        client.room.remove(client)
    }
})

handler.on(Opcodes.CREATE_ROOM, (data, client) => {
    const room = roomManager.findByOwner(client.id)

    if (room) {
        room.add(client)
    } else {
        roomManager.create().add(client)
    }
})

handler.on(Opcodes.JOIN_ROOM, 'string', (data, client) => {
    const room = roomManager.find(data)

    if (room) {
        room.add(client)
    } else {
        client.send({
            op: Opcodes.JOIN_ROOM,
            d: { status: 404 },
        })
    }
})

handler.on(Opcodes.LEAVE_ROOM, (data, client) => {
    if (client.room) {
        client.room.remove(client)
    }
})

handler.on(Opcodes.QUEUE_ADD, 'string', (data, client) => {
    if (!client.room) {
        return
    }

    client.room.queue.add({
        url: data,
        author: client.id,
    })
})

module.exports = handler
