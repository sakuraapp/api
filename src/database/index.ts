import { MikroORM } from '@mikro-orm/core'
import { TsMorphMetadataProvider } from '@mikro-orm/reflection'
import { BaseRepository } from './repositories/base.repository'
import { NameRepository } from './repositories/name.repository'
import { Discriminator } from './entities/discriminator.entity'
import { Name } from './entities/name.entity'
import { User, Credentials, Profile } from './entities/user.entity'
import { Service } from 'typedi'
import { logger } from '~/utils/logger'
import { Room } from './entities/room.entity'

@Service()
export default class Database {
    public orm: MikroORM

    // Repositories
    public name: NameRepository
    public user: BaseRepository<User>
    public room: BaseRepository<Room>

    async connect(): Promise<void> {
        this.orm = await MikroORM.init({
            entities: [User, Credentials, Profile, Name, Discriminator, Room],
            entityRepository: BaseRepository,
            dbName: 'sakura',
            type: 'mongo',
            host: process.env.DB_HOST || '127.0.0.1',
            port: Number(process.env.DB_PORT) || 27017,
            metadataProvider: TsMorphMetadataProvider,
            validate: true,
        })

        logger.info('Connected to database')
    }

    async init(): Promise<void> {
        await this.connect()

        this.name = this.orm.em.getRepository(Name)
        this.user = this.orm.em.getRepository(User)
        this.room = this.orm.em.getRepository(Room)
    }
}
