import { FastifyRequest, FastifyReply } from 'fastify'

export type NextFn = () => void
export type MiddlewareFn = (
    req: FastifyRequest,
    res: FastifyReply,
    next: NextFn
) => void
