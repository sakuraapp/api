const { Router } = require('express')
const router = Router()

router.get('/me', (req, res) => {
    res.json({
        id: req.user._id,
        discordId: req.user.credentials.userId,
        ...req.user.profile,
    })
})

module.exports = router
