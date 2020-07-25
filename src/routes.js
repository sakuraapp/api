const passport = require('passport')
const requireAuth = passport.authenticate('jwt', { session: false })

module.exports = (app) => {
    app.use('/auth', require('./controllers/auth'))
    app.use('/users', requireAuth, require('./controllers/user'))
    app.use('/rooms', requireAuth, require('./controllers/room'))
}
