import { ServerResponse } from "./types"

export async function getTopPostIds(scroll: number): Promise<ServerResponse> {
  let body = { Count: scroll * 10 }
  let res = await fetch(window.location.origin + "/api/get_top_shitposts", {
    method: "PUT",
    body: JSON.stringify(body)
  })
  return await res.json()
}

export async function getPosts(ids: number[]): Promise<ServerResponse> {
  let body = { ShitPostIds: ids }
  let res = await fetch(window.location.origin + "/api/get_saved_shitpost_list", {
    method: "PUT",
    body: JSON.stringify(body)
  })
  return await res.json()
}

export async function login(email: string, password: string): Promise<ServerResponse> {
  let body = { Login: email, Mdp: password }
  let res = await fetch(window.location.origin + "/api/login", {
    method: "PUT",
    body: JSON.stringify(body),
  })
  return await res.json()
}

export async function logout(): Promise<ServerResponse> {
  let res = await fetch(window.location.origin + "/api/logout")
  return await res.json()
}

export async function createAccount(email: string, password: string): Promise<ServerResponse> {
  let body = { Login: email, Mdp: password }
  let res = await fetch(window.location.origin + "/api/create_account", {
    method: "POST",
    body: JSON.stringify(body),
  })
  return await res.json()
}

export async function getPrivateProfile(): Promise<ServerResponse> {
  let res = await fetch(window.location.origin + "/api/get_private_profile")
  let json = await res.json()
  console.log(json)
  return json
}

export async function getRandomPost() {
  let res = await fetch("http://localhost:25565/api/random_shitpost")
  return await res.json()
}

export async function getSavedPost(id: number): Promise<ServerResponse> {
  let res = await fetch(window.location.origin + "/api/get_saved_shitpost", {})
  return await res.json()
}

export async function savePost(url: string, caption: string): Promise<ServerResponse> {
  url = url.replaceAll(" ", "_")
  let req = await fetch(window.location.origin + "/api/save_shitpost", {
    method: "POST",
    body: `{"url":"${url}", "caption": "${caption}"}`
  })
  return await req.json()
}