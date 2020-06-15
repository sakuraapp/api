const { Router } = require('express')
const querystring = require('querystring')

const passport = require('passport')
const jwt = require('jsonwebtoken')

const router = Router()

router.use((req, res, next) => {
    req.signToken = () => {
        if (req.user) {
            return jwt.sign({ id: req.user._id }, process.env.JWT_PRIVATE_KEY, {
                algorithm: 'RS256',
            })
        }

        return null
    }

    next()
})

router.get('/discord', (req, res) => {
    const opts = {
        response_type: 'code',
        client_id: process.env.DISCORD_CLIENT_ID,
        scope: process.env.DISCORD_SCOPES,
        redirect_uri: process.env.DISCORD_REDIRECT_URI,
    }

    res.end(`${process.env.DISCORD_OAUTH_URL}?${querystring.stringify(opts)}`)
})

router.get(
    '/login',
    passport.authenticate('discord', { session: false }),
    (req, res) => {
        res.status(200).json({ status: 200, token: req.signToken() })
    }
)

router.post(
    '/verify',
    passport.authenticate('jwt', { session: false }),
    (req, res) => {
        res.status(200).json({ status: 200 })
    }
)

module.exports = router
