const Schema = require('./schema')

const UserSchema = new Schema({
    credentials: {
        userId: String,
        accessToken: String,
        refreshToken: String,
    },
    profile: {
        username: String,
        discriminator: String,
        avatar: String,
    },
})

module.exports = UserSchema
