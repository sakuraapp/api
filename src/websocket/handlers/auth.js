const jwt = require('jsonwebtoken')

const Handler = require('./handler')
const handler = new Handler()

const Opcodes = require('@common/opcodes.json')

const User = require('~/models/user')
const logger = require('~/utils/logger')

function authenticate(token) {
    return new Promise((resolve, reject) => {
        jwt.verify(token, process.env.JWT_PUBLIC_KEY, (err, payload) => {
            if (err) reject(err)
            else {
                User.findOne({ _id: payload.id })
                    .then((user) => {
                        if (user) {
                            resolve(user)
                        } else {
                            reject(new Error("User doesn't exist."))
                        }
                    })
                    .catch(reject)
            }
        })
    })
}

handler.on(Opcodes.AUTHENTICATE, 'string', (token, client) => {
    authenticate(token)
        .then((user) => {
            client.user = user
            client.id = client.user.profile.id = user._id.toString()

            client.send({
                op: Opcodes.AUTHENTICATE,
                d: { socketId: client.socketId },
            })
        })
        .catch((err) => {
            logger.debug(err.stack || err)
            client.socket.close()
        })
})

handler.requireAuth = (data, client, next) => {
    if (client.user) {
        next()
    }
}

module.exports = handler
