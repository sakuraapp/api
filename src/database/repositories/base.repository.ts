import { FilterQuery, AnyEntity, EntityData } from '@mikro-orm/core'
import { EntityRepository } from '@mikro-orm/mongodb'

export class BaseRepository<T extends AnyEntity<T>> extends EntityRepository<
    T
> {
    public async findOrCreate(
        where: FilterQuery<T>,
        baseDoc?: EntityData<T>
    ): Promise<T> {
        if (!baseDoc) {
            baseDoc = {}
        }

        if (typeof where === 'object') {
            const conds = where as EntityData<T>

            baseDoc = {
                ...conds,
                ...baseDoc,
            }
        }

        const collection = this.em
            .getConnection()
            .getCollection(this.entityName.toString())

        const data = await collection.findOneAndUpdate(
            where as FilterQuery<unknown>,
            {
                $setOnInsert: baseDoc,
            },
            {
                upsert: true,
                returnOriginal: false,
            }
        )

        return this.map(data.value)
    }
}
