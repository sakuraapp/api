const fastify = require('fastify')
const passport = require('passport')
const passportHttp = require('passport/lib/http/request')

const got = require('got')
const cheerio = require('cheerio')
const URL = require('url')

const favicons = {
    'crunchyroll.com': 'https://www.crunchyroll.com/favicons/favicon-32x32.png',
}

exports.isDev = () => process.env.NODE_ENV !== 'production'

exports.createServer = () => {
    const opts = {
        logging: this.isDev(),
    }

    if (Number(process.env.USE_SSL)) {
        opts.https = {
            cert: process.env.SSL_CERT,
            key: process.env.SSL_KEY,
        }
    }

    const app = fastify(opts)

    // passport support
    app.decorateReply('setHeader', function (key, value) {
        this.header(key, value)
    })

    const senders = ['end', 'json']

    for (const sender of senders) {
        app.decorateReply(sender, function (data) {
            this.send(data)
        })
    }

    for (var i in passportHttp) {
        app.decorateRequest(i, passportHttp[i])
    }

    app.addHook('preHandler', passport.initialize())

    return app
}

exports.padLeft = function (str, len = 4) {
    return Array(len - String(str).length + 1).join('0') + str
}

exports.getSiteInfo = async (url) => {
    const domain = this.getDomain(url)

    let favicon
    let title

    try {
        const res = await got(url)
        const $ = cheerio.load(res.body)

        const getElement = (selector) => {
            const el = $(selector)

            return el.length > 0 ? el : null
        }

        const iconSelectors = [
            'link[rel="apple-touch-icon-precomposed"]',
            'link[rel="apple-touch-icon"]',
            'link[rel="shortcut icon"]',
            'link[rel="icon"]',
            'meta[itemprop="image"]@content',
        ]

        for (const selector of iconSelectors) {
            const args = selector.split('@')
            const attr = args[1] || 'href'

            const el = getElement(args[0])

            if (el) {
                favicon = el.attr(attr)

                if (favicon.startsWith('/')) {
                    favicon = res.url + favicon
                }

                break
            }

            title = $('title').text().trim()
        }
    } catch (err) {
        favicon = favicons[domain]
        title = domain
    }

    return {
        title,
        favicon,
    }
}

exports.getDomain = (url) => {
    const parsed = URL.parse(url)
    const { hostname } = parsed
    const parts = hostname.split('.')

    if (parts[0] === 'www') {
        parts.shift()
    }

    return parts.join('.')
}

exports.getYoutubeVideoId = (url) => {
    var regExp = /^.*(?:(?:youtu\.be\/|v\/|vi\/|u\/\w\/|embed\/)|(?:(?:watch)?\?v(?:i)?=|\&v(?:i)?=))([^#\&\?]+).*/
    var match = url.match(regExp)

    return match && match[1].length >= 11 ? match[1] : false
}
