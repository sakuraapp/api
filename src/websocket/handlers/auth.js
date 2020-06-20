const jwt = require('jsonwebtoken')

const Handler = require('./handler')
const handler = new Handler()

const User = require('~/models/user')
const logger = require('../../utils/logger')

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

handler.on('authenticate', 'string', (token, client) => {
    authenticate(token)
        .then((user) => {
            client.user = user
            client.send({ action: 'authenticated' })
        })
        .catch((err) => {
            logger.debug(err.stack || err)
            client.socket.close()
        })
})

module.exports = handler
