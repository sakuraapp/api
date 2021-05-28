import { OAuth2Strategy, IProfile } from './oauth2.strategy'
import { Strategy as PassportDiscordStrategy, Profile } from 'passport-discord'
import { StrategyToken } from './strategy.strategy'
import { Service } from 'typedi'
import querystring from 'qs'

@Service({ id: StrategyToken, multiple: true })
export class DiscordStrategy extends OAuth2Strategy<Profile> {
    public id = 'discord'
    public strategy = new PassportDiscordStrategy(
        {
            clientID: this.clientId,
            clientSecret: this.clientSecret,
            callbackURL: this.redirectUri,
            scope: this.scopes,
        },
        (accessToken, refreshToken, profile, cb) => {
            this.handleAuth(accessToken, refreshToken, profile)
                .then((user) => cb(null, user))
                .catch((err) => {
                    console.log(err.stack)
                    cb(err, null)
                })
        }
    )

    formatProfile(profile: Profile): IProfile {
        let { avatar } = profile

        if (avatar) {
            const avatarExt = avatar.startsWith('a_') ? 'gif' : 'png'

            avatar = `https://cdn.discordapp.com/avatars/${profile.id}/${avatar}.${avatarExt}`
        }

        return {
            id: profile.id,
            username: profile.username,
            avatar: avatar,
        }
    }

    getOAuth2Url(): string {
        const opts = {
            response_type: 'code',
            client_id: this.clientId,
            scope: this.scopesRaw,
            redirect_uri: this.redirectUri,
        }

        return `${this.oauthUrl}?${querystring.stringify(opts)}`
    }
}
