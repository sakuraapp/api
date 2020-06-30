const { EventEmitter } = require('events')
const Handler = require('./handlers/handler')

class MessageBroker extends EventEmitter {
    constructor() {
        super()

        /**
         * @type {Array<Handler>}
         */
        this.handlers = []
    }

    use(middleware, handler) {
        if (!handler) {
            handler = middleware
            middleware = null
        }

        if (middleware) handler.middleware.push(middleware)

        this.handlers.push(handler)
    }

    handle(packet, client) {
        if (!packet.action) return

        this.emit(packet.action, packet, client)

        for (const handler of this.handlers) {
            const listeners = handler.listeners.get(packet.action) || []

            for (const listener of listeners) {
                listener.handle(packet.data, client, packet.timestamp)
            }
        }
    }
}

module.exports = MessageBroker
