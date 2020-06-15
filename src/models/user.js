const { model } = require('mongoose')
const Schema = require('../schemas/user')

module.exports = model('user', Schema)
