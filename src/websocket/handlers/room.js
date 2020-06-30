const Handler = require('./handler')
const RoomManager = require('../room/manager')

const handler = new Handler()

handler.on('create room', (data, client) => {
    RoomManager.create().add(client)
})

module.exports = handler
