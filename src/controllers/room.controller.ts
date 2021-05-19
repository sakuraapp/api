import { FastifyReply, FastifyRequest } from 'fastify'
import { Controller, GET, POST } from 'fastify-decorators'
import { Container } from 'typedi'
import Database from '~/database'
import SessionController from './session.controller'

@Controller({ route: '/rooms' })
export default class RoomController extends SessionController {
    private database = Container.get(Database)

    @GET('/:roomId')
    async getRoom(request: FastifyRequest, reply: FastifyReply): Promise<void> {
        const { roomId } = request.params as Record<string, string>
        const room = await this.database.room.findOne({ id: roomId })

        if (room) {
            reply.send({
                room: {
                    id: room.id,
                    name: room.name,
                    owner: room.owner.profile.username,
                    private: room.private,
                },
            })
        } else {
            reply.status(404).send({})
        }
    }

    @POST('/:roomid/messages')
    sendMessage(request: FastifyRequest, reply: FastifyReply): void {
        console.log('loll')
    }
}
