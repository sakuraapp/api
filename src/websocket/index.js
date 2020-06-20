const { Server } = require('ws')

const logger = require('../utils/logger')
const Client = require('./client')

const MessageBroker = require('./messageBroker')
const handlers = require('./handlers')

module.exports = (appServer) => {
    const server = new Server({ server: appServer })
    const messageBroker = new MessageBroker()

    handlers(messageBroker)

    server.on('connection', (socket) => {
        logger.debug('A client connected.')

        new Client(socket, messageBroker)
    })
}
