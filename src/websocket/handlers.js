module.exports = (messageBroker) => {
    messageBroker.use(require('./handlers/auth'))
}
