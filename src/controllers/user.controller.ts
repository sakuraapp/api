import { FastifyReply, FastifyRequest } from 'fastify'
import { Controller, GET } from 'fastify-decorators'
import { Container } from 'typedi'
import Database from '~/database'
import { UserHelper } from '~/helpers/user.helper'
import { PassportService } from '~/services/passport.service'
import SessionController from './session.controller'

@Controller({ route: '/users' })
export default class UserController extends SessionController {
    private passport = Container.get(PassportService)
    private database = Container.get(Database)

    @GET('/@me')
    async getMyUser(
        request: FastifyRequest,
        reply: FastifyReply
    ): Promise<void> {
        const { user } = request

        if (user.credentials.providerId && user.credentials.accessToken) {
            const profile = await this.passport.getProfile(user)

            if (profile.avatar !== user.profile.avatar) {
                user.profile.avatar = profile.avatar

                await this.database.orm.em.flush()
            }
        }

        reply.send({
            user: UserHelper.build(user),
        })
    }
}
