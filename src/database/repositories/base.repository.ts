import { FilterQuery, AnyEntity, EntityData, Utils } from '@mikro-orm/core'
import { EntityRepository } from '@mikro-orm/mongodb'

export class BaseRepository<T extends AnyEntity<T>> extends EntityRepository<
    T
> {
    private exportEntity<A>(data: A): A {
        data = Object.assign({}, data)
        const meta = this.em.getMetadata().get(this.entityName.toString())

        Utils.renameKey(data, 'id', '_id')

        for (const i in data) {
            const prop = meta.properties[i]

            if (prop) {
                if (prop.fieldNames) {
                    Utils.renameKey(data, i, prop.fieldNames[0])
                }
            }
        }

        return data
    }

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
            this.exportEntity(where) as FilterQuery<unknown>,
            {
                $setOnInsert: this.exportEntity(baseDoc),
            },
            {
                upsert: true,
                returnOriginal: false,
            }
        )

        return this.map(data.value)
    }
}
