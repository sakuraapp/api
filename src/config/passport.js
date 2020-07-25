const passport = require('passport')
const passportJwt = require('passport-jwt')
const refresh = require('passport-oauth2-refresh')
const DiscordStrategy = require('passport-discord').Strategy
const JwtStrategy = passportJwt.Strategy
const { ExtractJwt } = passportJwt

const User = require('../models/user')
const Name = require('../models/name')

function formatUserProfile(profile, overrides, provider) {
    if (!provider && typeof overrides === 'string') {
        provider = overrides
        overrides = null
    }

    let res

    switch (provider) {
        case 'discord':
            let avatar

            if (profile.avatar) {
                const avatarExt = profile.avatar.startsWith('a_')
                    ? 'gif'
                    : 'png'

                avatar = `https://cdn.discordapp.com/avatars/${profile.id}/${profile.avatar}.${avatarExt}`
            }

            res = {
                username: profile.username,
                discriminator: null,
                avatar,
            }
    }

    if (overrides) {
        for (var i in overrides) {
            res[i] = overrides[i]
        }
    }

    return res
}

async function findOrCreateUser(accessToken, refreshToken, profile) {
    const name = profile.username

    var user = await User.findOne({ 'credentials.userId': String(profile.id) })

    if (!user) {
        const discriminator = await Name.findFreeDiscriminator(name)

        if (!discriminator) {
            throw new Error('No discriminators available for this username.')
        }

        user = await User.create({
            credentials: {
                userId: profile.id,
                accessToken,
                refreshToken,
            },
            profile: formatUserProfile(profile, { discriminator }, 'discord'),
        })

        await Name.addDiscriminator(name, discriminator, user._id)
    } else {
        if (refreshToken != user.credentials.refreshToken) {
            await User.updateOne(
                { _id: user._id },
                {
                    credentials: {
                        refreshToken: refreshToken,
                    },
                }
            )

            user.credentials.refreshToken = refreshToken
        }
    }

    return user
}

const discordStrategy = new DiscordStrategy(
    {
        clientID: process.env.DISCORD_CLIENT_ID,
        clientSecret: process.env.DISCORD_CLIENT_SECRET,
        callbackURL: process.env.DISCORD_REDIRECT_URI,
        scope: process.env.DISCORD_SCOPES.split(', '),
    },
    (accessToken, refreshToken, profile, cb) => {
        findOrCreateUser(accessToken, refreshToken, profile)
            .then((user) => cb(null, user))
            .catch((err) => cb(err))
    }
)

passport.use(discordStrategy)
refresh.use(discordStrategy)

passport.use(
    new JwtStrategy(
        {
            jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
            secretOrKey: process.env.JWT_PUBLIC_KEY,
        },
        (payload, cb) => {
            User.findOne({ _id: payload.id })
                .then((user) => {
                    if (user) {
                        cb(null, user)
                    } else {
                        cb(null, false)
                    }
                })
                .catch((err) => cb(err))
        }
    )
)

function attemptRefresh(user, strategyName) {
    return new Promise((resolve, reject) => {
        refresh.requestNewAccessToken(
            strategyName,
            user.credentials.refreshToken,
            (err, accessToken, refreshToken) => {
                if (err) {
                    reject(err)
                } else {
                    user.accessToken = accessToken
                    user.refreshToken = refreshToken

                    User.updateOne(
                        { _id: user._id },
                        {
                            credentials: {
                                accessToken,
                                refreshToken,
                            },
                        }
                    )
                        .then(resolve)
                        .catch(reject)
                }
            }
        )
    })
}

exports.fetchUserProfile = (user, strategyName) => {
    return new Promise((resolve, reject) => {
        const strategy = passport._strategies[strategyName]

        strategy.userProfile(user.credentials.accessToken, (err, res) => {
            if (err) {
                if (user.credentials.refreshToken) {
                    attemptRefresh(user, strategyName)
                        .then(() => {
                            exports
                                .fetchUserProfile(user, strategyName)
                                .then(resolve)
                                .catch(reject)
                        })
                        .catch(reject)
                } else {
                    reject(err)
                }
            } else {
                if (strategyName === 'discord') {
                    res = formatUserProfile(
                        res,
                        { discriminator: user.profile.discriminator },
                        'discord'
                    )
                }

                resolve(res)
            }
        })
    })
}

exports.requireAuth = passport.authenticate('jwt', { session: false })
