const { Router } = require('express')
const { fetchUserProfile } = require('../config/passport')
const User = require('../models/user')
const router = Router()

router.get('/me', async (req, res) => {
    const profile = await fetchUserProfile(req.user, 'discord')

    for (const key in profile) {
        req.user.profile[key] = profile[key]
    }

    res.json({
        id: req.user._id,
        discordId: req.user.credentials.userId,
        ...req.user.profile,
    })
})

module.exports = router
