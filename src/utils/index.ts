import { FastifyInstance } from 'fastify'
import passport from 'passport'
import passportHttp from 'passport/lib/http/request'
import { MiddlewareFn } from '~/middlewares/middleware.middleware'

export function padLeft(str: string, len = 4): string {
    return Array(len - str.length + 1).join('0') + str
}

// passport support for fastify
export function passportCompatiblity(server: FastifyInstance): void {
    server.decorateReply('setHeader', function (key: string, value: unknown) {
        this.header(key, value)
    })

    const senders = ['end', 'json']

    for (const sender of senders) {
        server.decorateReply(sender, function (data: unknown) {
            this.send(data)
        })
    }

    for (const i in passportHttp) {
        server.decorateRequest(i, passportHttp[i])
    }

    server.addHook(
        'preHandler',
        (passport.initialize() as unknown) as MiddlewareFn
    )
}
