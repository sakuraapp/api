const { requireAuth } = require('../config/passport')

module.exports = (fastify, opts, next) => {
    //fastify.addHook('preHandler', requireAuth)

    fastify.get('/', async (req, res) => {
        const response = await fastify.svcClient.api.get('/nodes')

        res.send(response.data)
    })

    next()
}
