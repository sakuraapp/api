const { Opcodes } = require('@sakuraapp/common')
const logger = require('~/utils/logger')
const utils = require('~/utils')

const DEFAULT_STATE = {
    playing: false,
    url: null,
    itemId: null,
    currentTime: null,
    playbackStart: null,
}

class Queue {
    constructor(room) {
        this.room = room

        this.items = []
        this.itemId = 0
        this.currentItem = null

        this.next()
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

    remove(id) {
        for (let i in this.items) {
            const item = this.items[i]

            if (item.id === id) {
                this.items.splice(i, 1)
                this.room.send({
                    op: Opcodes.QUEUE_REMOVE,
                    d: id,
                })

                return true
            }
        }

        return false
    }

    next(updateQueue = true) {
        const item = this.items.shift()

        if (item) {
            this.currentItem = item
            this.room.state = {
                playing: false,
                currentTime: 0,
                url: item.url,
                itemId: item.id,
            }

            if (updateQueue !== false) {
                this.room.send({
                    op: Opcodes.QUEUE_REMOVE,
                    d: item.id,
                })
            }
        } else {
            this.room.state = Object.assign({}, DEFAULT_STATE)
        }

        this.room.sendPlayerState()
    }
}

module.exports = Queue
