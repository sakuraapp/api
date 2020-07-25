const { Router } = require('express')
const { fetchUserProfile } = require('../config/passport')
const User = require('../models/user')
const router = Router()

router.get('/me', async (req, res) => {
    let profile
    let different

    try {
        profile = await fetchUserProfile(req.user, 'discord')
    } catch (err) {
        console.log(err)

        return res.status(401).end('401 Unauthorized')
    }

    for (const key in profile) {
        if (profile[key] !== req.user.profile[key]) {
            req.user.profile[key] = profile[key]
            different = true
        }
    }

    if (different) {
        await User.updateOne(
            { _id: req.user._id },
            {
                profile: req.user.profile.toObject(),
            }
        )
    }

    res.json({
        id: req.user._id,
        discordId: req.user.credentials.userId,
        ...req.user.profile,
    })
})

module.exports = router
