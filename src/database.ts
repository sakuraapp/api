import { Service } from 'typedi'
import { logger } from '~/utils/logger'
import {
    createDatabase,
    MikroORM,
    BaseRepository,
    NameRepository,
    Discriminator,
    Name,
    User,
    Credentials,
    Profile,
    Room,
} from '@sakuraapp/shared'

@Service()
export default class Database {
    public orm: MikroORM

    // Repositories
    public name: NameRepository
    public user: BaseRepository<User>
    public room: BaseRepository<Room>

    async connect(): Promise<void> {
        this.orm = await createDatabase({
            entities: [User, Credentials, Profile, Name, Discriminator, Room],
            dbName: 'sakura',
            host: process.env.DB_HOST || '127.0.0.1',
            port: Number(process.env.DB_PORT) || 27017,
            validate: true,
        })

        logger.info('Connected to the database')
    }

    async init(): Promise<void> {
        await this.connect()

        this.name = this.orm.em.getRepository(Name)
        this.user = this.orm.em.getRepository(User)
        this.room = this.orm.em.getRepository(Room)
    }
}
