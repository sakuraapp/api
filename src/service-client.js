// todo: move away from this god awful singleton structure
const roomManager = require('./websocket/room/manager')
const Room = require('./websocket/room/room')
const { Client } = require('@sakuraapp/service-client')
const client = new Client({ name: 'api' })

client.registerMethod('deploy', (packet) => {
    const { id, playingUrl } = packet.data.d

    if (!roomManager.rooms.has(id)) {
        roomManager.rooms.set(id, new Room(id, 2))
    }

    const room = roomManager.rooms.get(id)

    if (room) {
        room.queue.currentItem = {
            url: playingUrl,
            itemId: 0,
        }

        room.state = {
            ...room.queue.currentItem,
            playing: true,
            currentTime: 0,
        }

        room.sendPlayerState()
    }
})

module.exports = client
