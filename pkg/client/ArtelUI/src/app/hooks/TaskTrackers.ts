import {create} from 'zustand'

import {TaskTrackerInfo, AddTaskTrackerRequest, AddTaskTrackerResponse} from "@/app/api/artel/task_trackers.pb.ts"
import {taskTrackersService} from "@/processes/TaskTrackers.ts"

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
            const trackers = await taskTrackersService.list()
            set({trackers})
        } finally {
            set({loading: false})
        }
    },

    add: async (req: AddTaskTrackerRequest) => {
        const resp = await taskTrackersService.add(req)
        await get().fetch()
        return resp
    },

    remove: async (id: string) => {
        await taskTrackersService.remove(id)
        await get().fetch()
    },
}))
