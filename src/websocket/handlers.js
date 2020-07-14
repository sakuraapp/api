const authHandler = require('./handlers/auth')

module.exports = (messageBroker) => {
    messageBroker.use(authHandler)
    messageBroker.use(authHandler.requireAuth, require('./handlers/room'))
}
