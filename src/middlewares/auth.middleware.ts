import { FastifyRequest, FastifyReply } from 'fastify'
import { Container } from 'typedi'
import { StrategyManager } from '~/managers/strategy.manager'
import { IStrategy } from '~/strategies/strategy.strategy'
import { MiddlewareFn, NextFn } from './middleware.middleware'

export function authenticate({
    paramName,
    strategies,
}: {
    paramName: string
    strategies?: IStrategy[]
}): MiddlewareFn {
    return function (
        req: FastifyRequest,
        res: FastifyReply,
        next: NextFn
    ): void {
        if (!strategies) {
            strategies = Container.get(StrategyManager).strategies
        }

        const params = req.params as Record<string, string>
        const provider = params[paramName]
        const strategy = strategies.find(
            (strat) => strat.id === provider && strat.isOAuth2
        )

        if (strategy) {
            strategy.authenticate(req, res, next)
        }
    }
}
