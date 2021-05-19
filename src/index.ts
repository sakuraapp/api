import 'reflect-metadata'
import 'module-alias/register'

import dotenv from 'dotenv'
import Container from 'typedi'
import App from './app'

dotenv.config()

import './config/ssl'
import './config/jwt'

const app = Container.get(App)

app.init()
