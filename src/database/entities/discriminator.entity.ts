import { Entity, OneToOne, PrimaryKey, Property } from '@mikro-orm/core'
import { ObjectId } from '@mikro-orm/mongodb'
import { padLeft } from '~/utils'
import { User } from './user.entity'

export const MAX_DISCRIMINATOR_DIGITS = 4
export const MAX_DISCRIMINATOR_VALUE = Number(
    new Array(MAX_DISCRIMINATOR_DIGITS + 1).join('9')
)

export function toDiscriminator(num: number): string {
    return padLeft(num.toString(), MAX_DISCRIMINATOR_DIGITS)
}

@Entity({ collection: 'discriminators' })
export class Discriminator {
    @PrimaryKey()
    _id: ObjectId

    @Property()
    name: ObjectId

    @Property()
    value: string

    @OneToOne()
    ownerId: User
}
