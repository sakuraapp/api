const passport = require('passport')
const jwt = require('jsonwebtoken')

module.exports = (fastify, opts, next) => {
    fastify.decorateReply('signToken', function () {
        const req = this.request

        if (req.user) {
            return jwt.sign({ id: req.user._id }, process.env.JWT_PRIVATE_KEY, {
                algorithm: 'RS256',
            })
        }
    })

    fastify.get('/discord', (req, res) => {
        const opts = {
            response_type: 'code',
            client_id: process.env.DISCORD_CLIENT_ID,
            scope: process.env.DISCORD_SCOPES,
            redirect_uri: process.env.DISCORD_REDIRECT_URI,
        }

        const querystring = require('querystring')
        res.send(
            `${process.env.DISCORD_OAUTH_URL}?${querystring.stringify(opts)}`
        )
    })

    fastify.route({
        method: 'GET',
        url: '/login',
        preHandler: passport.authenticate('discord', { session: false }),
        handler(req, res) {
            res.status(200).send({ status: 200, token: res.signToken() })
        },
    })

    fastify.route({
        method: 'POST',
        url: '/verify',
        preHandler: passport.authenticate('jwt', { session: false }),
        handler(req, res) {
            res.status(200).send({ status: 200 })
        },
    })

    next()
}
