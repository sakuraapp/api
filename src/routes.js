const passport = require('passport')

module.exports = (app) => {
    app.use('/auth', require('./controllers/auth'))
    app.use(
        '/users',
        passport.authenticate('jwt', { session: false }),
        require('./controllers/user')
    )
}
