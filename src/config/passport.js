const passport = require('passport')
const passportJwt = require('passport-jwt')
const DiscordStrategy = require('passport-discord').Strategy
const JwtStrategy = passportJwt.Strategy
const { ExtractJwt } = passportJwt

const User = require('../models/user')
const Name = require('../models/name')

const OAUTH_MAX_AGE = 7200 // 2 hours

function makeProfile(profile, discriminator) {
    let avatar

    if (profile.avatar) {
        const avatarExt = profile.avatar.startsWith('a_') ? 'gif' : 'png'

        avatar = `https://cdn.discordapp.com/avatars/${profile.id}/${profile.avatar}.${avatarExt}`
    }

    return {
        username: profile.username,
        discriminator,
        avatar,
    }
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
            profile: makeProfile(profile, discriminator),
        })

        await Name.addDiscriminator(name, discriminator, user._id)
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

function getUserProfile(accessToken) {
    return new Promise((resolve, reject) => {
        discordStrategy.userProfile(accessToken, (err, profile) => {
            if (err) reject(err)
            else resolve(profile)
        })
    })
}

async function updateUserProfile(user) {
    const profile = await getUserProfile(user.credentials.accessToken)
    const newProfile = makeProfile(profile, user.profile.discriminator)
    let different = false

    for (var i in newProfile) {
        if (user[i] !== newProfile[i]) {
            different = true
            user[i] = newProfile[i]
        }
    }

    if (different) {
        await User.updateOne(
            { _id: user._id },
            { profile: newProfile, credentials: { lastQuery: Date.now() } }
        )
    }
}

passport.use(
    new JwtStrategy(
        {
            jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
            secretOrKey: process.env.JWT_PUBLIC_KEY,
        },
        (payload, cb) => {
            User.findOne({ _id: payload.id })
                .then(async (user) => {
                    if (user) {
                        const now = new Date().getTime()

                        if (
                            !user.credentials.lastQuery ||
                            now - user.credentials.lastQuery.getTime() >
                                OAUTH_MAX_AGE * 1000
                        ) {
                            try {
                                await updateUserProfile(user)
                            } catch (err) {}
                        }

                        cb(null, user)
                    } else {
                        cb(null, false)
                    }
                })
                .catch((err) => cb(err))
        }
    )
)
