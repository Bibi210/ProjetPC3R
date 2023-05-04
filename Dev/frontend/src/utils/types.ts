import React, { Dispatch, SetStateAction } from "react"

export type ServerResponse<T> = {
  Message: string,
  Success: boolean,
  Result: T
}

export type User = {
  Username: string
  LastSeen: string
  Posts: number[]
  Comments: number[]
  VotedComments: number[]
  VotedPosts: number[]
}

export type Post = {
  Id: number,
  Caption: string,
  Creator: string,
  Date: string,
  Upvotes: number,
  Url: string,
  CommentIds: number[]
}

export type Comment = {
  Id: number,
  Msg: {
    Content: string,
    Date: string,
    Sender: string
  },
  Upvotes: number
}

export type CurrentUserState = {
  get: User | null,
  set: React.Dispatch<SetStateAction<User | null>>,
  refresh: Function
}

export type PostComponentProps = {
  currentUserState: CurrentUserState
  loading: boolean,
  post: Post,
  setRefresh?: Dispatch<SetStateAction<boolean>>,
  randomMode: boolean,
}

export type CommentComponentProps = {
  currentUserState: CurrentUserState
  post: Post,
  showComments: boolean
}

export enum NotificationType {
  ERROR,
  INFO
}

export type Notification = {
  id: number,
  msg: string,
  type: NotificationType,
  show: boolean
}

export type SearchResults = {
  ShitPosts: number[],
  Users: string[]
}
