import {
    Entity,
    Property,
    PrimaryKey,
    Embedded,
    Embeddable,
    SerializedPrimaryKey,
    EntityRepositoryType,
} from '@mikro-orm/core'
import { ObjectId } from '@mikro-orm/mongodb'
import { BaseRepository } from '../repositories/base.repository'

@Embeddable()
export class Credentials {
    @Property({ nullable: true })
    userId?: string

    @Property({ nullable: true })
    providerId?: string

    @Property({ nullable: true })
    accessToken?: string

    @Property({ nullable: true })
    refreshToken?: string
}

@Embeddable()
export class Profile {
    @Property()
    username: string

    @Property()
    discriminator: string

    @Property({ nullable: true })
    avatar?: string
}

@Entity({ collection: 'users' })
export class User {
    [EntityRepositoryType]?: BaseRepository<User>

    @PrimaryKey()
    _id: ObjectId

    @SerializedPrimaryKey()
    id: string

    @Embedded({ object: true })
    credentials = new Credentials()

    @Embedded({ object: true })
    profile = new Profile()
}
