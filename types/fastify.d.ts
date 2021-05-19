import { User } from '~/database/entities/user.entity'

declare module 'fastify' {
    export interface FastifyRequest {
        user?: User
    }
}
