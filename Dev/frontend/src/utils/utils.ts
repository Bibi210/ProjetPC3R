import { upvotePost } from "./serverFunctions";
import { CurrentUserState, Post } from "./types";

export function isVoted(postId: number, currentUserState: CurrentUserState) {
  if (!currentUserState || !currentUserState.get) {
    // !notify
    return false
  } else if (currentUserState.get.VotedPosts) {
    return currentUserState.get.VotedPosts.includes(postId)
  } else {
    return false
  }
}

export function handleUpvotePost(post: Post, currentUserState: CurrentUserState, up: boolean) {
  upvotePost(post.Id, up ? 1 : 0).then(res => {
    if (res.Success && currentUserState.get) {
      let newPosts
      if (up)
        newPosts = currentUserState.get.VotedPosts
          ? [...currentUserState.get.VotedPosts, post.Id]
          : [post.Id]
      else
        newPosts = currentUserState.get.VotedPosts
          ? currentUserState.get.VotedPosts.filter(pId => pId != post.Id)
          : []

      currentUserState.set({
        Username: currentUserState.get.Username,
        Comments: currentUserState.get.Comments,
        Posts: currentUserState.get.Posts,
        LastSeen: currentUserState.get.LastSeen,
        VotedComments: currentUserState.get.VotedComments,
        VotedPosts: newPosts
      })
      currentUserState.refresh()
    }
  })
}
