const logger = require('../utils/logger')

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

            this.handleDisconnect()
            this.messageBroker.handle({ action: 'disconnect' }, this)
        })
    }

    handleDisconnect() {
        if (this.room) {
            this.room.remove(this)
        }
    }
}

module.exports = Client
