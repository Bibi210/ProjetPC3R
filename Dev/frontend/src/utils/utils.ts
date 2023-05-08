import { getComment, getPost, voteComment, votePost } from "./serverFunctions";
import { CurrentUserState, Post, Comment } from "./types";

export function isPostUpVoted(postId: number, currentUserState: CurrentUserState) {
  if (!currentUserState || !currentUserState.get) {
    return false
  } else if (currentUserState.get.UPVotedPosts) {
    return currentUserState.get.UPVotedPosts.includes(postId)
  } else {
    return false
  }
}

export function isPostDownVoted(postId: number, currentUserState: CurrentUserState) {
  if (!currentUserState || !currentUserState.get) {
    return false
  } else if (currentUserState.get.DOWNVotedPosts) {
    return currentUserState.get.DOWNVotedPosts.includes(postId)
  } else {
    return false
  }
}

export function isCommentUpVoted(commentId: number, currentUserState: CurrentUserState) {
  if (!currentUserState || !currentUserState.get) {
    return false
  } else if (currentUserState.get.UPVotedComments) {
    return currentUserState.get.UPVotedComments.includes(commentId)
  } else {
    return false
  }
}

export function isCommentDownVoted(commentId: number, currentUserState: CurrentUserState) {
  if (!currentUserState || !currentUserState.get) {
    return false
  } else if (currentUserState.get.DOWNVotedComments) {
    return currentUserState.get.DOWNVotedComments.includes(commentId)
  } else {
    return false
  }
}

export function handleUpVotePost(post: Post, currentUserState: CurrentUserState, up: boolean) {
  votePost(post.Id, up ? 1 : 0).then(res => {
    if (res.Success && currentUserState.get) {
      let newPosts
      if (up)
        newPosts = currentUserState.get.UPVotedPosts
          ? [...currentUserState.get.UPVotedPosts, post.Id]
          : [post.Id]
      else
        newPosts = currentUserState.get.UPVotedPosts
          ? currentUserState.get.UPVotedPosts.filter(pId => pId != post.Id)
          : []

      let downvotedPosts = currentUserState.get.DOWNVotedPosts
      if (downvotedPosts) {
        downvotedPosts = downvotedPosts.filter(postId => postId != post.Id)
      }
      currentUserState.set({
        Username: currentUserState.get.Username,
        Comments: currentUserState.get.Comments,
        Posts: currentUserState.get.Posts,
        LastSeen: currentUserState.get.LastSeen,
        UPVotedComments: currentUserState.get.UPVotedComments,
        UPVotedPosts: newPosts,
        DOWNVotedComments: currentUserState.get.DOWNVotedComments,
        DOWNVotedPosts: downvotedPosts
      })
    }
  })
}

export function handleVotePost(post: Post, setPost: Function | null, currentUserState: CurrentUserState, value: 1 | 0 | -1) {
  votePost(post.Id, value).then(res => {
    if (res.Success && currentUserState.get && setPost) {
      getPost(post.Id).then(postRes => {
        if (!postRes.Success) {
          alert(postRes.Message)
          console.error(postRes.Message)
        }
        setPost(postRes.Result)
        currentUserState.refresh()
      })
    }
  })
}

export function handleVoteComment(comment: Comment, comments: Comment[], setComments: Function, currentUserState: CurrentUserState) {
  if (!currentUserState || !currentUserState.get) {
    console.log("handleVoteComment:", "No connected user");
    return
  }

  let value: number
  if (currentUserState.get.UPVotedComments &&
    currentUserState.get.UPVotedComments.includes(comment.Id)) {
    value = -1
  } else if (currentUserState.get.DOWNVotedComments &&
    currentUserState.get.DOWNVotedComments.includes(comment.Id)) {
    value = 0
  } else {
    value = 1
  }

  voteComment(comment.Id, value).then((res) => {
    if (!res.Success) {
      alert(res.Message)
    } else {
      getComment(comment.Id).then(res => {
        if (!res.Success) {
          alert(res.Message)
        } else {
          let newComments: Comment[] = []
          for (const c of comments) {
            if (c.Id == comment.Id) {
              newComments.push(res.Result)
            } else {
              newComments.push(c)
            }
            setComments(newComments)
          }
          currentUserState.refresh()
        }
      })
    }
  })
}