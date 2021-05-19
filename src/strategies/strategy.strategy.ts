import { FastifyReply, FastifyRequest } from 'fastify'
import passport, { Strategy as PassportStrategy } from 'passport'
import { Inject, Service, Token } from 'typedi'
import Database from '~/database'
import { MiddlewareFn, NextFn } from '~/middlewares/middleware.middleware'

export interface IStrategy {
    id: string
    strategy: PassportStrategy
    isOAuth2: boolean
    authenticate: MiddlewareFn
}

export const StrategyToken = new Token<IStrategy>('strategies')

@Service()
export abstract class Strategy implements IStrategy {
    public abstract readonly id: string
    public abstract isOAuth2: boolean
    public abstract strategy: PassportStrategy

    @Inject()
    protected database: Database
    protected middleware: MiddlewareFn

    protected getConfig(name: string): string | null {
        return process.env[`${this.id.toUpperCase()}_${name.toUpperCase()}`]
    }

    public authenticate(
        req: FastifyRequest,
        res: FastifyReply,
        next: NextFn
    ): void {
        if (!this.middleware) {
            this.middleware = passport.authenticate(this.strategy.name, {
                session: false,
                failWithError: true,
            })
        }

        this.middleware(req, res, next)
    }
}
