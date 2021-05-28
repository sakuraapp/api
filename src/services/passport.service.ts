import LRU, { Options } from 'lru-cache'
import passport from 'passport'
import OAuth2Strategy from 'passport-oauth2'
import refresh from 'passport-oauth2-refresh'
import { Inject, Service } from 'typedi'
import { User } from '~/database/entities/user.entity'
import { StrategyManager } from '~/managers/strategy.manager'
import { IProfile } from '~/strategies/oauth2.strategy'
import { IStrategy } from '~/strategies/strategy.strategy'

const cacheOptions: Options<string, IProfile> = {
    max: 500, // 500 users
    maxAge: 10 * 60 * 1000, // 10 mins
}

@Service()
export class PassportService {
    private cache = new LRU<string, IProfile>(cacheOptions)

    @Inject()
    private strategyManager: StrategyManager

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

    async fetchProfile(user: User): Promise<IProfile> {
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

        const profile = await strat.userProfile(accessToken)

        this.cache.set(cacheKey, profile)

        return profile
    }
}
