const Opcodes = require('@common/opcodes.json')
const logger = require('~/utils/logger')
const utils = require('~/utils')

class Queue {
    constructor(room) {
        this.room = room

        this.items = []
        this.itemId = 0
        this.currentItem = null
    }

    add(item) {
        utils
            .getSiteInfo(item.url)
            .then(({ title, favicon }) => {
                item.title = title
                item.icon = favicon
                item.id = ++this.itemId

                this.items.push(item)

                if (this.items.length === 1 && !this.room.state.url) {
                    this.next(false)
                } else {
                    this.room.send({
                        op: Opcodes.QUEUE_ADD,
                        d: item,
                    })
                }
            })
            .catch(logger.error)
    }

    next(updateQueue = true) {
        const item = this.items.shift()

        this.currentItem = item
        this.room.state = {
            playing: false,
            currentTime: 0,
            url: item.url,
        }

        if (updateQueue !== false) {
            this.room.send({
                op: Opcodes.QUEUE_REMOVE,
                d: item.id,
            })
        }

        this.room.sendPlayerState()
    }
}

module.exports = Queue
