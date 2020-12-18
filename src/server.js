require('dotenv').config()
require('module-alias/register')

require('./config/ssl')
require('./config/jwt')
require('./config/passport')

const mongoose = require('mongoose')

const logger = require('./utils/logger')
const routes = require('./routes')
const websocket = require('./websocket')

const svcClient = require('./service-client')

const port = process.env.PORT
const version = process.env.npm_package_version

const { createServer } = require('./utils/index')

// todo: mongodb support for kubernetes
mongoose.Promise = Promise
mongoose.connect('mongodb://localhost/sakura', {
    useNewUrlParser: true,
    useUnifiedTopology: true,
})

const fastify = createServer()

fastify.decorate('svcClient', svcClient)
fastify.register(require('fastify-helmet'))
fastify.register(require('fastify-cors'))

svcClient.connect().then(() => {
    routes(fastify)
    websocket(fastify)

    console.log(`Sakura API v${version}\n`)

    fastify.listen(port, '0.0.0.0', () => {
        logger.write(`Listening on port ${port}`)
    })
})
