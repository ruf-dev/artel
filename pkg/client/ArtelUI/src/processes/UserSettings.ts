import useUser from "@/hooks/user/User.ts"
import {GetUserSettingsResponse, UserSettingsAPI} from "@/app/api/artel/user_settings.pb.ts"

export interface UserSettings {
    userPrompt: string
    likedOpenrouterModels: string[]
    lastUsedModel: string | undefined
}

// Shared by useLikedModels.ts and useLastUsedModel.ts: both hooks query
// GetUserSettings under this same key so they share one cached row instead of
// double-fetching when used together (see useSimpleChatController.ts).
export function userSettingsQueryKey() {
    return ["user-settings"] as const
}

function toUserSettings(res: GetUserSettingsResponse): UserSettings {
    return {
        userPrompt: res.userPrompt ?? "",
        likedOpenrouterModels: res.likedOpenrouterModels ?? [],
        lastUsedModel: res.lastUsedModel,
    }
}

export interface IUserSettingsService {
    getUserSettings: () => Promise<UserSettings>
    setLikedModels: (modelIds: string[]) => Promise<void>
    setLastUsedModel: (model: string) => Promise<void>
}

// Thin wrapper over the generated user_settings.pb.ts client, reshaping proto
// responses into plain TS objects — mirrors processes/SimpleChat.ts's shape.
export class UserSettingsService implements IUserSettingsService {
    async getUserSettings(): Promise<UserSettings> {
        const res = await UserSettingsAPI.GetUserSettings({}, useUser.getState().auth.getInitReq())
        return toUserSettings(res)
    }

    async setLikedModels(modelIds: string[]): Promise<void> {
        await UserSettingsAPI.SetLikedModels({likedOpenrouterModels: modelIds}, useUser.getState().auth.getInitReq())
    }

    async setLastUsedModel(model: string): Promise<void> {
        await UserSettingsAPI.SetLastUsedModel({model}, useUser.getState().auth.getInitReq())
    }
}

export const userSettingsService = new UserSettingsService()
