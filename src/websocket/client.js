const logger = require('../utils/logger')
const Opcodes = require('@common/opcodes.json')

class Client {
    constructor(socket, messageBroker) {
        this.socket = socket
        this.messageBroker = messageBroker

        this.bindEvents()
    }

    get username() {
        return this.user.profile.username
    }

    get profile() {
        if (!this.user) return null

        return {
            id: this.id,
            ...this.user.profile.toObject(),
        }
    }

    send(data) {
        if (!this.socket) return

        if (!data.t) {
            data.t = new Date().getTime()
        }

        const packet = JSON.stringify(data)

        this.socket.send(packet)
        logger.debug(`Sent: ${packet}`)
    }

    bindEvents() {
        this.socket.on('message', (data) => {
            logger.debug(`Received: ${data}`)

            try {
                const packet = JSON.parse(data)

                this.messageBroker.handle(packet, this)
            } catch (err) {
                logger.error(err.stack || err)
            }
        })

        this.socket.on('close', () => {
            logger.debug('A client has disconnected')

            this.socket = null
            this.messageBroker.handle({ op: Opcodes.DISCONNECT }, this)
        })
    }
}

module.exports = Client
