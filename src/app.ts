import { FastifyInstance } from 'fastify'
import { IncomingMessage, Server, ServerResponse } from 'http'
import helmet from 'fastify-helmet'
import cors from 'fastify-cors'
import { bootstrap } from 'fastify-decorators'
import { Service, Inject } from 'typedi'
import Database from './database'
import { ServiceManager } from './managers/service.manager'
import { logger } from './utils/logger'
import { createServer, passportCompatiblity } from './utils'
import AuthController from './controllers/auth.controller'
import UserController from './controllers/user.controller'
import RoomController from './controllers/room.controller'

@Service()
export default class App {
    public readonly address = '0.0.0.0'
    public readonly port = Number(process.env.PORT || '3000')

    public server: FastifyInstance<Server, IncomingMessage, ServerResponse>

    @Inject()
    public database: Database

    @Inject()
    public serviceManager: ServiceManager

    listen(): Promise<void> {
        return new Promise((resolve, reject) => {
            this.server.listen(this.port, this.address, (err) => {
                if (err) {
                    reject(err)
                } else {
                    resolve()
                }
            })
        })
    }

    async init(): Promise<void> {
        this.serviceManager.init()
        await this.database.init()

        this.server = createServer()

        this.server.register(helmet)
        this.server.register(cors)

        this.server.register(bootstrap, {
            controllers: [AuthController, UserController, RoomController],
            prefix: '/api/v1',
        })

        this.server.setErrorHandler((err, request, reply) => {
            console.log(err.stack)
            reply.send({
                message: err.message,
                code: err.code,
            })
        })

        this.server.addHook(
            'preSerialization',
            (request, reply, payload, done) => {
                if (typeof payload === 'object') {
                    done(null, {
                        status: reply.statusCode,
                        ...payload,
                    })
                } else {
                    done(null, payload)
                }
            }
        )

        passportCompatiblity(this.server)

        await this.listen()

        logger.info(`Sakura API is listening on port ${this.port}`)
    }
}
