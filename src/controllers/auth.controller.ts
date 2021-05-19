import { FastifyRequest, FastifyReply } from 'fastify'
import { Controller, GET } from 'fastify-decorators'
import { Container } from 'typedi'
import jwt from 'jsonwebtoken'
import { authenticate } from '~/middlewares/auth.middleware'
import { StrategyManager } from '~/managers/strategy.manager'

@Controller({ route: '/oauth2' })
export default class AuthController {
    private strategyManager = Container.get(StrategyManager)

    signToken(request: FastifyRequest): string | null {
        if (request.user) {
            return jwt.sign(
                { id: request.user._id },
                process.env.JWT_PRIVATE_KEY,
                { algorithm: 'RS256' }
            )
        }
    }

    @GET('/:providerId')
    getOAuthUrl(request: FastifyRequest, reply: FastifyReply): void {
        const { providerId } = request.params as Record<string, string>
        const strategy = this.strategyManager.getOAuth2(providerId)

        if (strategy) {
            reply.send(strategy.getOAuth2Url())
        } else {
            reply.status(404).send({})
        }
    }

    @GET('/:providerId/login', {
        preHandler: authenticate({ paramName: 'providerId' }),
    })
    loginOAuth2(request: FastifyRequest, reply: FastifyReply): void {
        reply.send({ token: this.signToken(request) })
    }
}
