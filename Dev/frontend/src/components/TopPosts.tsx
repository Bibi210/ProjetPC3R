import { useEffect, useState } from 'react'
import Post from './Post'
import { CurrentUserState, Post as PostType } from '../utils/types'
import { getPosts, getTopPostIds } from '../utils/serverFunctions'
import { Button, CircularProgress, Grid } from '@mui/material'

function TopPosts({
  currentUserState,
}: {
  currentUserState: CurrentUserState
}) {
  const [postLimit, setPostLimit] = useState(0)
  const [posts, setPosts] = useState<PostType[]>([])
  const [loading, setLoading] = useState(true)

  function addMorePosts() {
    let oldPostLimit = postLimit
    setPostLimit(postLimit + 1)
    setLoading(true)
    // setPostLimit doesn't update postLimit until a new render has happened
    getTopPostIds(oldPostLimit + 1).then((idsRes) => {
      getPosts(idsRes.Result).then((postsRes) => {
        setLoading(false)
        if (postsRes.Success && postsRes.Result) {
          setPosts(postsRes.Result)
        }
      })
    })
  }

  useEffect(() => {
    addMorePosts()
  }, [])

  return (
    <>
      {posts.map((p) => (
        <Post
          key={p.Url + p.Creator + p.Date}
          loading={false}
          post={p}
          randomMode={false}
          currentUserState={currentUserState}
        />
      ))}
      <Grid container justifyContent='center' marginBottom={4}>
        <Grid item>
          <Button
            variant='contained'
            onClick={() => addMorePosts()}
            style={{
              width: '840px',
            }}
          >
            {loading ? <CircularProgress /> : 'Load more posts'}
          </Button>
        </Grid>
      </Grid>
    </>
  )
}

export default TopPosts
