const fs = require('fs')
const colors = require('colors')

class Logger {
    constructor({ savePath = null, levels, loop = false }) {
        this.path = savePath
        this.saveCache = Boolean(savePath)

        this.levels = []
        this.cache = []

        this.createLevels(levels)

        if (loop !== false)
            this.interval = setInterval(this.saveLoop.bind(this), 5 * 60 * 1000)
    }

    /**
     * @private
     */
    saveLoop() {
        return new Promise((resolve, reject) => {
            if (this.cache.length === 0) return

            fs.appendFile(
                this.path,
                `${this.cache.join('\n')}\n`,
                'utf8',
                (err) => {
                    if (err) {
                        console.warn('Failed to save logs. WTF!')
                        console.error(err)

                        reject(err)
                    } else {
                        this.cache = []

                        this.log(null, 'Saved logs successfully.', {
                            save: false,
                        })
                        resolve()
                    }
                }
            )
        })
    }

    /**
     * Adds a log message to be saved
     * @private
     * @param {string} data
     */
    save(data) {
        if (!this.saveCache) return

        data = `[${this.logDate()}]${data}`

        this.cache.push(data)
    }

    /**
     * Logs a message.
     * @param {string} level
     * @param {string} data - Message to log
     * @param {object} [opts]
     */
    log(level, data, opts = {}) {
        var msg
        var raw

        var log = true
        var save = false

        var levelObj = typeof level === 'object' ? level : this.findLevel(level)

        if (!levelObj)
            levelObj = this.findLevel('info') || {
                name: 'info',
                color: 'green',
            }

        const { color } = levelObj
        const name = levelObj.name.toUpperCase()
        msg = levelObj.key || `[${name}]`

        if (opts.prefix) msg = `[${opts.prefix}]${msg}`
        if (opts.suffix) msg += `[${opts.suffix}]`

        msg += ' > '
        raw = msg

        if (color && colors[color]) msg = colors[color](msg)
        if (opts.color && colors[opts.color]) msg = colors[opts.color](msg)

        msg += data
        raw += data

        if (levelObj.log === false) log = false
        if (opts.log === true) log = true
        else if (opts.log === false)
            // Need it to be explicitly false, that's why I use else if instead of else
            log = false

        if (levelObj.save === false) save = false
        if (opts.save === true) save = true
        else if (opts.save === false)
            // ^
            save = false

        if (log) console.log(msg)

        if (save) this.save(raw)

        if (levelObj.fatal || opts.fatal) {
            if (opts.alert !== false) this.log('warning', 'Shutting down...')

            process.exit()
        }
    }

    logDate(d) {
        function padLeft(num, base, chr) {
            var len = String(base || 10).length - String(num).length + 1

            return len > 0 ? new Array(len).join(chr || '0') + num : num
        }

        if (!d) d = new Date()

        return `${[
            padLeft(d.getMonth() + 1),
            padLeft(d.getDate()),
            d.getFullYear(),
        ].join('-')} ${[
            padLeft(d.getHours()),
            padLeft(d.getMinutes()),
            padLeft(d.getSeconds()),
        ].join(':')}`
    }

    createLevel(level) {
        const name = level.name.toLowerCase()
        var funcName = name

        if (level.function) funcName = level.function

        this[funcName] = (...args) => {
            this.log(level, ...args)
        }
        level.function = { f: this[funcName], name: funcName }
        this.levels.push(level)
    }

    createLevels(levels) {
        for (const i in levels) {
            this.createLevel(levels[i])
        }
    }

    createSubLevel(level) {
        level.parent = level.parent.toLowerCase()
        const parent = this.levels.find(
            (lev) => lev.name.toLowerCase() === level.parent
        )

        if (!parent) {
            this.log(
                'warning',
                `Parent ${level.parent} not found for sublevel ${level.name}.`
            )

            return
        }

        if (!level.color) level.color = parent.color

        level.key = `[${parent.name.toUpperCase()}][${level.name.toUpperCase()}]`
        level.sublevel = true

        this.createLevel(level)
    }

    findLevel(name) {
        return this.levels.find((level) => level.name === name.toLowerCase())
    }

    removeLevel(level) {
        const func = level.function

        this[func.name] = null
        delete this[func.name]

        for (var i in this.levels) {
            if (this.levels[i] === level) {
                this.levels.splice(i, 1)

                break
            }
        }
    }

    removeLevels(removeSubs) {
        for (var level of this.levels) {
            if (!level.sublevel || removeSubs === true) this.removeLevel(level)
        }
    }
}

module.exports = new Logger({
    levels: [
        { name: 'info', function: 'write', color: 'green', save: false },
        { name: 'warning', function: 'warn', color: 'yellow' },
        { name: 'error', color: 'red' },
        { name: 'fatal', color: 'red', fatal: true },
    ],
})
