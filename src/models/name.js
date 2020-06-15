const { model } = require('mongoose')
const utils = require('../utils')
const Schema = require('../schemas/name')

const MAX_DISCRIMINATOR_VALUE = 9999
const getDiscrimByValue = (value, discrims) =>
    discrims.find((discrim) => discrim.value === value)

class DiscriminatorError extends Error {}

class NameModel {
    static async fetch(name) {
        const res = await this.findOrCreate({
            value: name,
        })

        return res.doc
    }
    static async findDiscriminators(name) {
        const res = await this.fetch(name)

        return res.discriminators
    }

    static async findFreeDiscriminator(name) {
        const discriminators = await this.findDiscriminators(name)

        for (var i = 1; i < MAX_DISCRIMINATOR_VALUE + 1; i++) {
            const discrim = utils.padLeft(i, 4)

            if (!getDiscrimByValue(discrim, discriminators)) {
                return discrim
            }
        }

        return null
    }

    static async addDiscriminator(name, value, ownerId) {
        const nameObj = await this.fetch(name)
        const { discriminators } = nameObj

        if (getDiscrimByValue(value, discriminators)) {
            throw new DiscriminatorError(
                `A discriminator with the same value already exists for name "${name}"`
            )
        }

        return this.updateOne(
            { _id: nameObj._id },
            {
                $push: {
                    discriminators: { value, ownerId },
                },
            }
        )
    }
}

Schema.loadClass(NameModel)

module.exports = model('name', Schema)
