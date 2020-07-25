const { Router } = require('express')
const roomManager = require('~/websocket/room/manager')
const router = Router()

router.get('/:roomId', (req, res) => {
    const { roomId } = req.params
    const room = roomManager.find(roomId)

    if (room) {
        res.json(room.getPublicInfo()).end()
    } else {
        res.status(404).end()
    }
})

router.post('/:roomId/messages', (req, res) => {
    const { roomId } = req.params
    const room = roomManager.find(roomId)

    if (room) {
        const client = room.getClientById(req.user._id.toString())

        if (client) {
            const id = room.sendMessage(req.body, client)

            return res.status(200).json({ id })
        }

        return res.status(403).end()
    }

    res.status(404).end()
})

module.exports = router
