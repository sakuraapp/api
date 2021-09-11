import LRU, { Options } from 'lru-cache'
import passport from 'passport'
import OAuth2Strategy from 'passport-oauth2'
import refresh from 'passport-oauth2-refresh'
import { Inject, Service } from 'typedi'
import { User } from '@sakuraapp/shared'
import { StrategyManager } from '~/managers/strategy.manager'
import { IOAuth2Strategy, IProfile } from '~/strategies/oauth2.strategy'
import { IStrategy } from '~/strategies/strategy.strategy'
import Database from '~/database'

const cacheOptions: Options<string, IProfile> = {
    max: 500, // 500 users
    maxAge: 10 * 60 * 1000, // 10 mins
}

@Service()
export class PassportService {
    private cache = new LRU<string, IProfile>(cacheOptions)

    @Inject()
    private strategyManager: StrategyManager

    @Inject()
    private database: Database

    get strategies(): IStrategy[] {
        return this.strategyManager.strategies
    }

    init(): void {
        const { strategies } = this

        for (const strat of strategies) {
            passport.use(strat.strategy)

            if (strat.isOAuth2) {
                refresh.use(strat.strategy as OAuth2Strategy)
            }
        }
    }

    public refresh(user: User, strat: IOAuth2Strategy): Promise<void> {
        return new Promise((resolve, reject) => {
            refresh.requestNewAccessToken(
                strat.strategy.name,
                user.credentials.refreshToken,
                (err, accessToken, refreshToken) => {
                    if (err) {
                        reject(err)
                    } else {
                        user.credentials.accessToken = accessToken
                        user.credentials.refreshToken = refreshToken

                        this.database.orm.em.flush().then(resolve).catch(reject)
                    }
                }
            )
        })
    }

    async fetchProfile(user: User, strat: IOAuth2Strategy): Promise<IProfile> {
        const { accessToken, refreshToken } = user.credentials

        let profile: IProfile

        try {
            profile = await strat.userProfile(accessToken)
        } catch (err) {
            if (refreshToken) {
                await this.refresh(user, strat)

                profile = await this.fetchProfile(user, strat)
            } else {
                throw err
            }
        }

        return profile
    }

    async getProfile(user: User): Promise<IProfile> {
        const { userId, providerId, accessToken } = user.credentials
        const strat = this.strategyManager.getOAuth2(providerId)

        if (!strat) {
            throw new Error(`Provider doesn't exist`)
        }

        const cacheKey = `${providerId}.${userId}`
        const cached = this.cache.get(cacheKey)

        if (cached) {
            return cached
        }

        const profile = await this.fetchProfile(user, strat)

        this.cache.set(cacheKey, profile)

        return profile
    }
}
