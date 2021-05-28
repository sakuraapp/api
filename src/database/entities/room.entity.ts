import {
    Entity,
    EntityRepositoryType,
    Enum,
    OneToOne,
    PrimaryKey,
    Property,
} from '@mikro-orm/core'
import { BaseRepository } from '../repositories/base.repository'
import { User } from './user.entity'

export enum RoomType {
    VideoSync = 1,
    VirtualBrowser = 2,
}

@Entity({ collection: 'rooms' })
export class Room {
    [EntityRepositoryType]?: BaseRepository<Room>

    @PrimaryKey()
    _id: string

    @Property()
    name: string

    @OneToOne({ fieldName: 'ownerId' })
    owner: User

    @Property()
    private: boolean

    @Enum(() => RoomType)
    type: RoomType
}
