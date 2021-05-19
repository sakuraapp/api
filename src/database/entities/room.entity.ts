import { Entity, Enum, OneToOne, PrimaryKey, Property } from '@mikro-orm/core'
import { ObjectId } from '@mikro-orm/mongodb'
import { User } from './user.entity'

export enum RoomType {
    VideoSync = 1,
    VirtualBrowser = 2,
}

@Entity({ collection: 'rooms' })
export class Room {
    @PrimaryKey()
    _id: ObjectId

    @Property()
    id: string

    @Property()
    name: string

    @OneToOne({ fieldName: 'ownerId' })
    owner: User

    @Property()
    private: boolean

    @Enum(() => RoomType)
    type: RoomType
}
