const logger = require('../utils/logger')
const Client = require('./client')

const MessageBroker = require('./messageBroker')
const handlers = require('./handlers')

module.exports = (fastify) => {
    const messageBroker = new MessageBroker()

    handlers(messageBroker)

    fastify.register(require('fastify-websocket'), {
        handle(conn) {
            logger.debug('A client connected.')

            new Client(conn.socket, messageBroker)
        },
        options: {},
    })
}
