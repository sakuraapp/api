import { FastifyReply, FastifyRequest } from 'fastify'
import { Controller, GET, POST } from 'fastify-decorators'
import { nanoid } from 'nanoid'
import { Container } from 'typedi'
import Database from '~/database'
import SessionController from './session.controller'

@Controller({ route: '/rooms' })
export default class RoomController extends SessionController {
    private database = Container.get(Database)

    @GET('/:roomId')
    async getRoom(request: FastifyRequest, reply: FastifyReply): Promise<void> {
        const { roomId } = request.params as Record<string, string>
        const room = await this.database.room.findOne({ _id: roomId })

        if (room) {
            reply.send({
                room: {
                    id: room._id,
                    name: room.name,
                    owner: room.owner.profile.username,
                    private: room.private,
                },
            })
        } else {
            reply.status(404).send({})
        }
    }

    @POST('/')
    async createRoom(
        request: FastifyRequest,
        reply: FastifyReply
    ): Promise<void> {
        const { user } = request
        const room = await this.database.room.findOrCreate(
            { owner: user._id },
            {
                id: nanoid(), // todo: maybe not generate nanoid by default and check if room exists first?
                name: `${user.profile.username}'s room`, // todo: localization ????
            }
        )

        reply.send({
            id: room._id,
        })
    }

    @POST('/:roomid/messages')
    sendMessage(request: FastifyRequest, reply: FastifyReply): void {
        console.log('loll')
    }
}
