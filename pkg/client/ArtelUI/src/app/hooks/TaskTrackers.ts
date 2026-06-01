import {create} from 'zustand'

import {TaskTrackerInfo, TaskTrackersAPI, AddTaskTrackerRequest, AddTaskTrackerResponse} from "@/app/api/artel/task_trackers.pb.ts"
import useUser from "@/hooks/user/User.ts"

interface TaskTrackersState {
    trackers: TaskTrackerInfo[]
    loading: boolean
    fetch: () => Promise<void>
    add: (req: AddTaskTrackerRequest) => Promise<AddTaskTrackerResponse>
    remove: (id: string) => Promise<void>
}

export const useTaskTrackers = create<TaskTrackersState>((set, get) => ({
    trackers: [],
    loading: false,

    fetch: async () => {
        set({loading: true})
        try {
            const {auth} = useUser.getState()
            const resp = await TaskTrackersAPI.ListTaskTrackers({}, auth.getInitReq())
            set({trackers: resp.trackers ?? []})
        } finally {
            set({loading: false})
        }
    },

    add: async (req: AddTaskTrackerRequest) => {
        const {auth} = useUser.getState()
        const resp = await TaskTrackersAPI.AddTaskTracker(req, auth.getInitReq())
        await get().fetch()
        return resp
    },

    remove: async (id: string) => {
        const {auth} = useUser.getState()
        await TaskTrackersAPI.DeleteTaskTracker({id}, auth.getInitReq())
        await get().fetch()
    },
}))
