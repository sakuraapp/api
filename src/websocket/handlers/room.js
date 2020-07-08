const Handler = require('./handler')
const Opcodes = require('@common/opcodes.json')
const RoomManager = require('../room/manager')

const handler = new Handler()

handler.on(Opcodes.CREATE_ROOM, (data, client) => {
    RoomManager.create().add(client)
})

module.exports = handler
