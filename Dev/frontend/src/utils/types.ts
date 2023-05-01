import { Dispatch, SetStateAction } from "react"
import { getPrivateProfile } from "./serverFunctions";

export type ServerResponse = {
  Message: string,
  Success: boolean,
  Result: any
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
  Msg: {
    Content: string,
    Date: string,
    Sender: string
  },
  Upvotes: number
}

export type PostComponentProps = {
  loading: boolean,
  post: Post,
  setRefresh?: Dispatch<SetStateAction<boolean>>,
  randomMode: boolean,
  showCommentBtn: boolean
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

let currentUser: User | null = null

export function getCurrentUser()  {
  if (!currentUser) {
    getPrivateProfile().then(res => {
      if (res.Success) {
        currentUser = res.Result
        return currentUser
      } else {
        return null
      }
    })
  } else {
    return currentUser
  }
}