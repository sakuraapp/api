import {
    Discriminator,
    MAX_DISCRIMINATOR_VALUE,
    toDiscriminator,
} from '../entities/discriminator.entity'
import { Name } from '../entities/name.entity'
import { BaseRepository } from './base.repository'

export class NameRepository extends BaseRepository<Name> {
    public async findFreeDiscriminator(
        name: string
    ): Promise<Discriminator | null> {
        const doc = await this.findOrCreate(
            { value: name },
            { discriminators: [] }
        )

        await doc.discriminators.init()

        for (let i = 0; i < MAX_DISCRIMINATOR_VALUE; i++) {
            if (!doc.discriminators[i]) {
                const discrim = new Discriminator()

                this.em.persist(discrim)

                discrim.name = doc._id
                discrim.value = toDiscriminator(i + 1)

                doc.discriminators.add(discrim)

                return discrim
            }

            if (!doc.discriminators[i].ownerId) {
                return doc.discriminators[i]
            }
        }
    }
}
