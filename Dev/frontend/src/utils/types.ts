import { Dispatch, SetStateAction } from "react"

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
  Caption: string,
  Creator: string,
  Date: string,
  Upvotes: number,
  Url: string
}

export type PostComponentProps = {
  loading: boolean,
  caption: string,
  src: string,
  setRefresh?: Dispatch<SetStateAction<boolean>>,
  random: boolean,
  comments: boolean
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