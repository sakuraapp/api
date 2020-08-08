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

handler.on(Opcodes.PLAYER_STATE, (data, client) => {
    if (!client.room || !data) {
        return
    }

    const { room } = client

    if (!room.hasPermission('VIDEO_REMOTE', client)) {
        return
    }

    const now = new Date().getTime()
    const latency = now - data.t

    const { playing, currentTime } = data

    let valid = false

    if (!isNaN(currentTime) && currentTime !== null) {
        room.state.currentTime = currentTime + latency / 1000
        room.state.playbackStart = now
        valid = true
    }

    if (typeof playing === 'boolean') {
        room.setPlaying(playing, false)
        valid = true
    }

    if (valid) room.sendPlayerState()
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

    if (!client.hasPermission('QUEUE_ADD')) {
        return
    }

    client.room.queue.add({
        url: data,
        author: client.id,
    })
})

module.exports = handler
