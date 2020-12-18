const Handler = require('./handler')
const { Opcodes, RoomTypes } = require('@sakuraapp/common')
const roomManager = require('../room/manager')
const { getDomain, getYoutubeVideoId } = require('~/utils')
const svcClient = require('~/service-client')

const handler = new Handler()

handler.on(Opcodes.DISCONNECT, (data, client) => {
    if (client.room) {
        client.room.remove(client)
    }
})

// data is room type
handler.on(Opcodes.CREATE_ROOM, 'number', (data, client) => {
    const room = roomManager.findByOwner(client.id)

    if (room) {
        room.add(client)
    } else {
        const room = roomManager.create(data)

        room.add(client)

        if (data === 2) {
            svcClient.call('room/deploy', {
                id: room.id,
                location: 'all',
                priority: false,
            })
        }
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
    let latency = now - data.t

    const { playing, currentTime } = data

    if (!playing) {
        latency = 0
    }

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

    const domain = getDomain(data)

    switch (domain) {
        case 'youtube.com':
            const id = getYoutubeVideoId(data)

            if (id) {
                data = `https://www.youtube.com/embed/${id}`
            }
    }

    client.room.queue.add({
        url: data,
        author: client.id,
    })
})

handler.on(Opcodes.QUEUE_REMOVE, 'number', (data, client) => {
    if (!client.room) {
        return
    }

    if (!client.hasPermission('QUEUE_EDIT')) {
        return
    }

    client.room.queue.remove(data)
})

handler.on(Opcodes.VIDEO_END, 'number', (data, client) => {
    if (!client.room) {
        return
    }

    const { room } = client
    const clients = room.findClientsWithPermissions(['VIDEO_REMOTE'])

    if (clients.length === 0 || clients.includes(client)) {
        if (data === room.queue.itemId) {
            room.queue.next(true)
        }
    }
})

handler.on(Opcodes.VIDEO_SKIP, (data, client) => {
    if (!client.room) {
        return
    }

    if (!client.hasPermission('VIDEO_REMOTE')) {
        return
    }

    client.room.queue.next(true)
})

handler.on(Opcodes.SET_ROOM_TYPE, 'number', (data, client) => {
    if (client.room.ownerId !== client.id) {
        return // only owner can change room type
    }

    if (!RoomTypes[data]) {
        return
    }

    client.room.type = data
})

handler.on(Opcodes.DISPATCH_CONTROL, (data, client) => {
    if (!client.room) {
        return
    }

    if (!client.hasPermission('VIDEO_REMOTE')) {
        return
    }

    svcClient
        .call('controller/dispatch', {
            ...data,
            room: client.room.id,
        })
        .catch(console.error)
})

module.exports = handler
