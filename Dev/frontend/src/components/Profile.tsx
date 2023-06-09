import {
  Avatar,
  Button,
  Card,
  CardContent,
  CardHeader,
  CircularProgress,
  Container,
  Grid,
  Typography,
} from '@mui/material'
import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import Post from './Post'
import { CurrentUserState, Post as PostType } from '../utils/types'
import { deleteAccount, getPosts } from '../utils/serverFunctions'

function Profile({ currentUserState }: { currentUserState: CurrentUserState }) {
  const [loading, setLoading] = useState(false)
  const [posts, setPosts] = useState<PostType[]>([])

  useEffect(() => {
    if (
      currentUserState &&
      currentUserState.get &&
      currentUserState.get.Posts
    ) {
      getPosts(currentUserState.get.Posts).then((postsRes) => {
        if (postsRes.Success) {
          setPosts(postsRes.Result ? postsRes.Result : [])
        } else {
          console.log(postsRes.Message)
          alert(postsRes.Message)
        }
        setLoading(false)
      })
    }
  }, [])
  return (
    <Container>
      <Card color='white' variant='outlined' >
        <CardHeader
          avatar={
            loading ? (
              <CircularProgress />
            ) : currentUserState && currentUserState.get ? (
              <Avatar> {initials(currentUserState.get.Username)}</Avatar>
            ) : (
              <Avatar></Avatar>
            )
          }
          title={
            loading ? (
              ''
            ) : currentUserState && currentUserState.get ? (
              <Typography variant='h6'>
                {currentUserState.get.Username}
              </Typography>
            ) : (
              'No connected user'
            )
          }
          action={
            loading ? (
              <CircularProgress />
            ) : (
              currentUserState &&
              currentUserState.get && (
                <Grid container spacing={1}>
                  <Grid item>
                    <Button
                      variant='contained'
                      style={{ backgroundColor: '#EF5350', color: 'white' }}
                      onClick={() => {
                        if (
                          confirm(
                            'Are you sure you want to delete your account?'
                          )
                        ) {
                          deleteAccount()
                        }
                      }}
                    >
                      Delete account
                    </Button>
                  </Grid>
                  <Grid item>
                    <Link to='/logout'>
                      <Button variant='contained'>Logout</Button>
                    </Link>
                  </Grid>
                </Grid>
              )
            )
          }
        ></CardHeader>

        <CardContent>
          {loading ? (
            <CircularProgress />
          ) : currentUserState && currentUserState.get ? (
            posts.length == 0 ? (
              <Grid container justifyContent='center'>
                <Grid item>
                  <Typography variant='h5'>Nothing here yet</Typography>
                </Grid>
              </Grid>
            ) : (
              <Grid container justifyContent='center' marginBottom={4}>
                {posts.map((post: PostType) => (
                  <Grid item key={post.Id}>
                    <Post
                      currentUserState={currentUserState}
                      loading={false}
                      post={post}
                      randomMode={false}
                    />
                  </Grid>
                ))}
              </Grid>
            )
          ) : (
            <Link to='/login'>
              <Button
                fullWidth
                style={{ backgroundColor: '#3F51B5', color: 'white' }}
              >
                Login
              </Button>
            </Link>
          )}
        </CardContent>
      </Card>
    </Container>
  )
}

export function initials(name: string): string {
  return name
    .split(' ')
    .map((n) => n[0].toUpperCase())
    .join()
}

export default Profile
