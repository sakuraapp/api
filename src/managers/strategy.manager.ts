import { Container, Service, InjectMany } from 'typedi'
import { IStrategy, StrategyToken } from '~/strategies/strategy.strategy'
import { DiscordStrategy } from '~/strategies/discord.strategy'
import { JwtStrategy } from '~/strategies/jwt.strategy'
import { IOAuth2Strategy } from '~/strategies/oauth2.strategy'

Container.import([DiscordStrategy, JwtStrategy])

@Service()
export class StrategyManager {
    @InjectMany(StrategyToken)
    public strategies: IStrategy[]

    get<T extends IStrategy>(id: string): T | null {
        return this.strategies.find((strat) => strat.id === id) as T
    }

    getOAuth2<T extends IOAuth2Strategy>(id: string): T | null {
        const strat = this.get<T>(id)

        if (strat) {
            return strat.isOAuth2 ? strat : null
        }
    }
}
