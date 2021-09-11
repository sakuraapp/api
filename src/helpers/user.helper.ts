import { User, Profile } from '@sakuraapp/shared'

export interface UserInfo extends Profile {
    id: string
}

export class UserHelper {
    static build(user: User): UserInfo {
        return {
            id: user.id,
            ...user.profile,
        }
    }
}
