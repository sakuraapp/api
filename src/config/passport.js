const passport = require('passport')
const passportJwt = require('passport-jwt')
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
    }

    return user
}

passport.use(
    new DiscordStrategy(
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
)

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

exports.fetchUserProfile = (user, strategyName) => {
    return new Promise((resolve, reject) => {
        const strategy = passport._strategies[strategyName]

        strategy.userProfile(user.credentials.accessToken, (err, res) => {
            if (err) reject(err)
            else {
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
