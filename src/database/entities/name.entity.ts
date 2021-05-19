import {
    Property,
    Entity,
    PrimaryKey,
    Collection,
    EntityRepositoryType,
    ManyToMany,
} from '@mikro-orm/core'
import { ObjectId } from '@mikro-orm/mongodb'
import { NameRepository } from '../repositories/name.repository'
import { Discriminator } from './discriminator.entity'

@Entity({
    collection: 'names',
    customRepository: () => NameRepository,
})
export class Name {
    [EntityRepositoryType]?: NameRepository

    @PrimaryKey()
    _id: ObjectId

    @Property()
    value: string

    @ManyToMany(() => Discriminator)
    discriminators = new Collection<Discriminator>(this)
}
