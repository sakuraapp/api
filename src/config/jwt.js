const fs = require('fs')
const path = require('path')

if (!process.env.JWT_PUBLIC_KEY) {
    process.env.JWT_PUBLIC_KEY = fs.readFileSync(
        path.join(__dirname, '/public.pem')
    )
}

if (!process.env.JWT_PRIVATE_KEY) {
    process.env.JWT_PRIVATE_KEY = fs.readFileSync(
        path.join(__dirname, 'private.pem')
    )
}
