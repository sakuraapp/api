import { FastifyReply, FastifyRequest } from 'fastify'
import { Hook } from 'fastify-decorators'
import { Container } from 'typedi'
import { StrategyManager } from '~/managers/strategy.manager'
import { NextFn } from '~/middlewares/middleware.middleware'

export default class SessionController {
    private strategyManager = Container.get(StrategyManager)

    @Hook('preHandler')
    auth(request: FastifyRequest, reply: FastifyReply, next: NextFn): void {
        this.strategyManager.get('jwt').authenticate(request, reply, next)
    }
}
