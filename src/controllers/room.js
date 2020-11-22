const { Router } = require('express')
const bodyParser = require('body-parser')
const roomManager = require('~/websocket/room/manager')
const router = Router()

router.use(bodyParser.json())

router.get('/:roomId', (req, res) => {
    const { roomId } = req.params
    const room = roomManager.find(roomId)

    if (room) {
        res.json(room.getPublicInfo()).end()
    } else {
        res.status(404).end()
    }
})

router.post('/:roomId/messages', (req, res, next) => {
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

            return res.status(200).json({ id })
        }

        return res.status(403).end()
    }

    res.status(404).end()
})

module.exports = router
