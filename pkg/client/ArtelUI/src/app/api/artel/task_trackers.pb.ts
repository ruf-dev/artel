/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export type TaskTrackerInfo = {
  id?: string;
  type?: string;
  name?: string;
  createdAt?: string;
};

export type TrelloBoardInfo = {
  id?: string;
  name?: string;
};

export type AddTaskTrackerRequest = {
  type?: string;
  apiKey?: string;
  apiToken?: string;
};

export type AddTaskTrackerResponse = {
  tracker?: TaskTrackerInfo;
  boards?: TrelloBoardInfo[];
};

export type AddTaskTracker = Record<string, never>;

export type ListTaskTrackersRequest = Record<string, never>;

export type ListTaskTrackersResponse = {
  trackers?: TaskTrackerInfo[];
};

export type ListTaskTrackers = Record<string, never>;

export type DeleteTaskTrackerRequest = {
  id?: string;
};

export type DeleteTaskTrackerResponse = Record<string, never>;

export type DeleteTaskTracker = Record<string, never>;

export type ListTrelloBoardsRequest = {
  trackerId?: string;
};

export type ListTrelloBoardsResponse = {
  boards?: TrelloBoardInfo[];
};

export type ListTrelloBoards = Record<string, never>;

export class TaskTrackersAPI {
  static AddTaskTracker(this:void, req: AddTaskTrackerRequest, initReq?: fm.InitReq): Promise<AddTaskTrackerResponse> {
    return fm.fetchRequest<AddTaskTrackerResponse>(`/api/task-trackers/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListTaskTrackers(this:void, req: ListTaskTrackersRequest, initReq?: fm.InitReq): Promise<ListTaskTrackersResponse> {
    return fm.fetchRequest<ListTaskTrackersResponse>(`/api/task-trackers/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DeleteTaskTracker(this:void, req: DeleteTaskTrackerRequest, initReq?: fm.InitReq): Promise<DeleteTaskTrackerResponse> {
    return fm.fetchRequest<DeleteTaskTrackerResponse>(`/api/task-trackers/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListTrelloBoards(this:void, req: ListTrelloBoardsRequest, initReq?: fm.InitReq): Promise<ListTrelloBoardsResponse> {
    return fm.fetchRequest<ListTrelloBoardsResponse>(`/api/task-trackers/trello/boards`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}