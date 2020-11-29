module.exports = (fastify) => {
    fastify.register(require('./controllers/auth'), { prefix: '/auth' })
    fastify.register(require('./controllers/user'), { prefix: '/users' })
    fastify.register(require('./controllers/room'), { prefix: '/rooms' })
}
