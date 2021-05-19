import { readFileSync } from 'fs'
import { join } from 'path'

if (Number(process.env.USE_SSL)) {
    process.env.SSL_CERT = readFileSync(join(__dirname, 'cert.pem'), 'utf-8')
    process.env.SSL_KEY = readFileSync(join(__dirname, 'key.pem'), 'utf-8')
}
