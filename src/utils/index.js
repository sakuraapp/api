const http = require('http')
const https = require('https')

exports.createServer = function (app) {
    return Number(process.env.USE_SSL)
        ? https.createServer(
              {
                  cert: process.env.SSL_CERT,
                  key: process.env.SSL_KEY,
              },
              app
          )
        : http.createServer(app)
}

exports.padLeft = function (str, len = 4) {
    return Array(len - String(str).length + 1).join('0') + str
}
