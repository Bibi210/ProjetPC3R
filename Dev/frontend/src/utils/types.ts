import { Dispatch, SetStateAction } from "react"

export type User = {
  Username: string
  Posts: number[]
}

export type ServerResponse = {
  Message: string,
  Success: boolean,
  Result: any
}

export type Post = {
  Caption: string,
  Creator: string,
  Date: string,
  Upvotes: number,
  Url: string
}

export type PostProps = {
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