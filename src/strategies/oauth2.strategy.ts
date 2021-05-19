import PassportOAuth2Strategy from 'passport-oauth2'
import { IStrategy, Strategy } from './strategy.strategy'
import { User } from '~/database/entities/user.entity'

// Abstract Profile
export interface IProfile {
    id: string
    username: string
    avatar?: string
}

export interface IOAuth2Strategy extends IStrategy {
    getOAuth2Url?(): string
    userProfile(accessToken: string): Promise<IProfile>
}

export abstract class OAuth2Strategy<T = unknown> extends Strategy {
    public isOAuth2 = true
    public abstract strategy: PassportOAuth2Strategy
    public abstract formatProfile?(profile: T): IProfile
    public abstract getOAuth2Url(): string

    public async handleAuth(
        accessToken: string,
        refreshToken: string,
        sourceProfile: T
    ): Promise<User> {
        const profile = this.formatProfile(sourceProfile)

        const name = profile.username
        const credentials = {
            userId: profile.id,
            providerId: this.id,
        }

        let user = await this.database.user.findOne({ credentials })

        if (!user) {
            const discrim = await this.database.name.findFreeDiscriminator(name)

            if (!discrim) {
                throw new Error('Too many users are using this username')
            }

            user = new User()

            discrim.ownerId = user

            user.credentials = credentials
            user.profile = {
                username: profile.username,
                discriminator: discrim.value,
            }

            this.database.orm.em.persist(user)
        }

        user.credentials.accessToken = accessToken
        user.credentials.refreshToken = refreshToken

        user.profile.avatar = profile.avatar

        await this.database.orm.em.flush()

        return user
    }

    public async userProfile(accessToken: string): Promise<IProfile> {
        return new Promise((resolve, reject) => {
            this.strategy.userProfile(accessToken, (err, profile: T) => {
                if (err) {
                    reject(err)
                } else {
                    resolve(this.formatProfile(profile))
                }
            })
        })
    }

    get oauthUrl(): string {
        return this.getConfig('oauth_url')
    }

    get clientId(): string {
        return this.getConfig('client_id')
    }

    get clientSecret(): string {
        return this.getConfig('client_secret')
    }

    get redirectUri(): string {
        return this.getConfig('redirect_uri')
    }

    get scopesRaw(): string | null {
        return this.getConfig('scopes')
    }

    get scopes(): string[] {
        return this.scopesRaw?.split(', ') || []
    }
}
