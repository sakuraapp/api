class Listener {
    constructor(fn, typeDefs, handler) {
        /**
         * @type {Function}
         */
        this.fn = fn

        /**
         * @type {object}
         */
        this.typeDefs = typeDefs

        /**
         * @type {Handler}
         */
        this.handler = handler
    }

    checkTypes(data, defs) {
        if (!defs) return true

        switch (typeof defs) {
            case 'string':
                return typeof data === defs
                break
            case 'object':
                for (var i in data) {
                    if (!this.checkTypes(data[i], defs[i])) {
                        return false
                    }
                }

                return true
                break
            case 'function':
                return data.constructor === defs || data instanceof defs
        }
    }

    handle(data, client) {
        this.handler.runMiddleware([data, client], () => {
            if (this.typeDefs) {
                if (!this.checkTypes(data, this.typeDefs)) {
                    return
                }
            }

            this.fn(data, client)
        })
    }
}

class Handler {
    constructor() {
        /**
         * @prop {Map<String, Listener>}
         */
        this.listeners = new Map()

        /**
         * @prop {Array<Function>}
         */
        this.middleware = []
    }

    /**
     * @method
     * @param {string} action
     * @param {string|Function|object} typeDefinitions
     * @param {Function} callback
     */
    /**
     * @method
     * @param {string} action
     * @param {Function} callback
     */
    on(action, typeDefs, cb) {
        if (!cb) {
            cb = typeDefs
            typeDefs = undefined
        }

        if (!this.listeners.has(action)) {
            this.listeners.set(action, [])
        }

        this.listeners.get(action).push(new Listener(cb, typeDefs, this))
    }

    async runMiddleware(args, cb) {
        for (let fn of this.middleware) {
            const res = new Promise((resolve) => {
                fn(...args, resolve)
            })

            await res
        }

        cb()
    }
}

module.exports = Handler
