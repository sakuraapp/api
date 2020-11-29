require('dotenv').config()
require('module-alias/register')

require('./config/ssl')
require('./config/jwt')
require('./config/passport')

const mongoose = require('mongoose')

const logger = require('./utils/logger')
const routes = require('./routes')
const websocket = require('./websocket')

const port = process.env.API_PORT
const version = process.env.npm_package_version

const { createServer } = require('./utils/index')

mongoose.Promise = Promise
mongoose.connect('mongodb://localhost/sakura', {
    useNewUrlParser: true,
    useUnifiedTopology: true,
})

const fastify = createServer()

fastify.register(require('fastify-helmet'))
fastify.register(require('fastify-cors'))

routes(fastify)
websocket(fastify)

console.log(`Sakura API v${version}\n`)

fastify.listen(port, () => {
    logger.write(`Listening on port ${port}`)
})
