require('dotenv').config()
require('module-alias/register')

require('./config/ssl')
require('./config/jwt')
require('./config/passport')

const express = require('express')

const helmet = require('helmet')
const cors = require('cors')

const mongoose = require('mongoose')
const passport = require('passport')

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

const app = express()
const server = createServer(app)

app.use(helmet())
app.use(cors())
app.use(passport.initialize())

routes(app)
websocket(server)

console.log(`Sakura API v${version}\n`)

server.listen(port, () => {
    logger.write(`Listening on port ${port}`)
})
