const Schema = require('./schema')
const { ObjectId } = require('mongoose')

const NameSchema = new Schema({
    value: String,
    discriminators: [
        {
            value: String,
            ownerId: ObjectId,
        },
    ],
})

module.exports = NameSchema
