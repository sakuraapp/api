const fs = require('fs')
const path = require('path')

if (Number(process.env.USE_SSL)) {
    process.env.SSL_CERT = fs.readFileSync(
        path.join(__dirname, 'cert.pem'),
        'utf-8'
    )

    process.env.SSL_KEY = fs.readFileSync(
        path.join(__dirname, 'key.pem'),
        'utf-8'
    )
}
