import { Room } from '@sakuraapp/shared'
import { UserHelper, UserInfo } from './user.helper'

export interface RoomInfo {
    id: string
    name: string
    owner: UserInfo
    private: boolean
}

export class RoomHelper {
    static build(room: Room): RoomInfo {
        return {
            id: room._id,
            name: room.name,
            owner: UserHelper.build(room.owner),
            private: room.private,
        }
    }
}
