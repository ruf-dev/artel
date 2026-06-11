import useUser from "@/hooks/user/User.ts"
import {TaskTrackerInfo, TaskTrackersAPI, AddTaskTrackerRequest, AddTaskTrackerResponse} from "@/app/api/artel/task_trackers.pb.ts"

export interface ITaskTrackersService {
    list: () => Promise<TaskTrackerInfo[]>
    add: (req: AddTaskTrackerRequest) => Promise<AddTaskTrackerResponse>
    remove: (id: string) => Promise<void>
}

export class TaskTrackersService implements ITaskTrackersService {
    async list(): Promise<TaskTrackerInfo[]> {
        const res = await TaskTrackersAPI.ListTaskTrackers({}, useUser.getState().auth.getInitReq())
        return res.trackers ?? []
    }

    async add(req: AddTaskTrackerRequest): Promise<AddTaskTrackerResponse> {
        return TaskTrackersAPI.AddTaskTracker(req, useUser.getState().auth.getInitReq())
    }

    async remove(id: string): Promise<void> {
        await TaskTrackersAPI.DeleteTaskTracker({id}, useUser.getState().auth.getInitReq())
    }
}

export const taskTrackersService = new TaskTrackersService()
