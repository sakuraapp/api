import { Inject, Service } from 'typedi'
import { PassportService } from '~/services/passport.service'

@Service()
export class ServiceManager {
    @Inject()
    public passport: PassportService

    init(): void {
        this.passport.init()
    }
}
