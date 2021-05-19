import { Service } from 'typedi'
import { Strategy, StrategyToken } from './strategy.strategy'
import { ExtractJwt, Strategy as PassportJwtStrategy } from 'passport-jwt'

@Service({ id: StrategyToken, multiple: true })
export class JwtStrategy extends Strategy {
    public id = 'jwt'
    public isOAuth2 = false
    public strategy = new PassportJwtStrategy(
        {
            jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
            secretOrKey: this.getConfig('public_key'),
        },
        (payload, cb) => {
            this.database.user
                .findOne(payload.id)
                .then((user) => {
                    if (user) {
                        cb(null, user)
                    } else {
                        cb(null, false)
                    }
                })
                .catch((err) => cb(err))
        }
    )
}
