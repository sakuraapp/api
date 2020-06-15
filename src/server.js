require('dotenv').config()
require('./config/jwt')
require('./config/passport')

const express = require('express')

const helmet = require('helmet')
const cors = require('cors')

const mongoose = require('mongoose')
const passport = require('passport')

const logger = require('./utils/logger')
const routes = require('./routes')

const port = process.env.API_PORT
const { version } = require('../package')

mongoose.Promise = Promise
mongoose.connect('mongodb://localhost/sakura', {
    useNewUrlParser: true,
    useUnifiedTopology: true,
})

const app = express()

app.use(helmet())
app.use(cors())
app.use(passport.initialize())

routes(app)

console.log(`Sakura API v${version}\n`)

app.listen(port, () => {
    logger.write(`Listening on port ${port}`)
})
