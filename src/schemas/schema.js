const mongoose = require('mongoose')
const findOrCreate = require('mongoose-findorcreate')

class Schema extends mongoose.Schema {
    constructor(...args) {
        super(...args)

        this.plugin(findOrCreate)
    }
}

module.exports = Schema
