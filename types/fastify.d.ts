import { User } from '@sakuraapp/shared'

declare module 'fastify' {
    export interface FastifyRequest {
        user?: User
    }
}
