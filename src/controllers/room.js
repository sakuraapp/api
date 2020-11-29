const { requireAuth } = require('~/config/passport')
const roomManager = require('~/websocket/room/manager')

module.exports = (fastify, opts, next) => {
    fastify.addHook('preHandler', requireAuth)

    fastify.get('/:roomId', (req, res) => {
        const { roomId } = req.params
        const room = roomManager.find(roomId)

        if (room) {
            res.send(room.getPublicInfo())
        } else {
            res.status(404)
        }
    })

    fastify.post('/:roomId/messages', (req, res, next) => {
        const socketId = req.headers['x-socket-id']

        if (!socketId) {
            return next()
        }

        const { roomId } = req.params
        const room = roomManager.find(roomId)

        if (room) {
            const client = room.getClientBySocketId(socketId)

            if (client) {
                const id = room.sendMessage(req.body.content, client)

                return res.status(200).send({ id })
            }

            return res.status(403)
        }

        res.status(404)
    })

    next()
}
