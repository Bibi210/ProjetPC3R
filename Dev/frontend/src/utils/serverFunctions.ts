import { Comment, Post, SearchResults, ServerResponse, User } from "./types"

export async function getTopPostIds(scroll: number): Promise<ServerResponse<number[]>> {
  let body = { Count: scroll * 10 }
  let res = await fetch(window.location.origin + "/api/get_top_shitposts", {
    method: "PUT",
    body: JSON.stringify(body)
  })
  return await res.json()
}

export async function getPosts(ids: number[]): Promise<ServerResponse<Post[]>> {
  let body = { ShitPostIds: ids }
  let res = await fetch(window.location.origin + "/api/get_saved_shitpost_list", {
    method: "PUT",
    body: JSON.stringify(body)
  })
  return await res.json()
}

export async function login(email: string, password: string): Promise<ServerResponse<string | null>> {
  let body = { Login: email, Mdp: password }
  let res = await fetch(window.location.origin + "/api/login", {
    method: "PUT",
    body: JSON.stringify(body),
  })
  return await res.json()
}

export async function logout(): Promise<ServerResponse<null>> {
  let res = await fetch(window.location.origin + "/api/logout")
  return await res.json()
}

export async function createAccount(login: string, password: string): Promise<ServerResponse<null>> {
  let body = { Login: login, Mdp: password }
  let res = await fetch(window.location.origin + "/api/create_account", {
    method: "POST",
    body: JSON.stringify(body),
  })
  return await res.json()
}

export async function getPrivateProfile(): Promise<ServerResponse<User | null>> {
  let res = await fetch(window.location.origin + "/api/get_private_profile")
  return await res.json()
}

export async function getPublicProfile(username: string): Promise<ServerResponse<User | null>> {
  let body = { Username: username }
  let res = await fetch(window.location.origin + "/api/get_public_profile", {
    method: "PUT",
    body: JSON.stringify(body)
  })
  return await res.json()
}

export async function getRandomPost(): Promise<ServerResponse<any>> {
  let res = await fetch(window.location.origin + "/api/random_shitpost")
  return await res.json()
}

export async function savePost(url: string, caption: string): Promise<ServerResponse<null>> {
  url = url.replaceAll(" ", "_")
  let res = await fetch(window.location.origin + "/api/save_shitpost", {
    method: "POST",
    body: `{"url":"${url}", "caption": "${caption}"}`
  })
  return await res.json()
}

export async function search(query: string): Promise<ServerResponse<SearchResults>> {
  let body = { Query: query }
  let res = await fetch(window.location.origin + "/api/search", {
    method: "PUT",
    body: JSON.stringify(body)
  })
  return await res.json()
}

export async function getComments(ids: number[]): Promise<ServerResponse<Comment[]>> {
  let body = { CommentIds: ids }
  let res = await fetch(window.location.origin + "/api/get_comment_list", {
    method: "PUT",
    body: JSON.stringify(body)
  })
  return await res.json()
}

export async function saveComment(postId: number, content: string)
  : Promise<ServerResponse<{ Id: number } | null>> {
  let body = { ShitPostId: postId, Content: content }
  let res = await fetch(window.location.origin + "/api/post_comment", {
    method: "POST",
    body: JSON.stringify(body)
  })
  return await res.json()
}

export async function upvotePost(postId: number, value: number): Promise<ServerResponse<null>> {
  let body = { ShitPostId: postId, Value: value }
  let res = await fetch(window.location.origin + "/api/post_shitpost_vote", {
    method: "POST",
    body: JSON.stringify(body)
  })
  return await res.json()
}

export async function upvoteComment(commentId: number, value: number): Promise<ServerResponse<any>> {
  let body = { CommentId: commentId, Value: value }
  let res = await fetch(window.location.origin + "/api/post_comment_vote", {
    method: "POST",
    body: JSON.stringify(body)
  })
  return await res.json()
}
