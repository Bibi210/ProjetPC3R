import {
  Avatar,
  Box,
  Card,
  CardActionArea,
  CardContent,
  CardMedia,
  Collapse,
  Container,
  IconButton,
  InputAdornment,
  TextField,
  Typography,
} from '@mui/material'
import { useEffect, useRef, useState } from 'react'
import { getPosts, getPublicProfile, search } from '../utils/serverFunctions'
import {
  CurrentUserState,
  Post,
  SearchResults,
  ServerResponse,
  User,
} from '../utils/types'
import { initials } from './Profile'
import { filetype } from './Post'
import { KeyboardArrowDown, KeyboardArrowUp } from '@mui/icons-material'
import { handleVotePost, isPostDownVoted, isPostUpVoted } from '../utils/utils'

function Search({ currentUserState }: { currentUserState: CurrentUserState }) {
  const [searchResults, setSearchResults] = useState<SearchResults>({
    ShitPosts: [],
    Users: [],
  })
  let keyTimeout = useRef(0)

  function onChange(value: string) {
    if (keyTimeout.current != 0) {
      clearTimeout(keyTimeout.current)
    }
    keyTimeout.current = setTimeout(() => {
      if (value != '') {
        search(value).then((searchRes) => {
          if (searchRes.Success) {
            let result: SearchResults = searchRes.Result
            if (result.ShitPosts == null) {
              result.ShitPosts = []
            }
            if (result.Users == null) {
              result.Users = []
            }
            setSearchResults(result)
          }
        })
      } else {
        setSearchResults({ ShitPosts: [], Users: [] })
      }
    }, 500)
  }

  return (
    <Container>
      <TextField
        fullWidth={true}
        onChange={(e) => onChange(e.target.value)}
        InputProps={{
          startAdornment: (
            <InputAdornment position='start'>🔍 search</InputAdornment>
          ),
        }}
        style={{ marginBottom: '20px' }}
      />
      <SearchResultUsers usernames={searchResults.Users} />
      <SearchResultPosts
        currentUserState={currentUserState}
        postIds={searchResults.ShitPosts}
      />
    </Container>
  )
}

function SearchResultUsers({ usernames }: { usernames: string[] }) {
  const [users, setUsers] = useState<User[]>([])
  useEffect(() => {
    let userRequests: Promise<ServerResponse<User | null>>[] = []
    for (const username of usernames) {
      userRequests.push(getPublicProfile(username))
    }
    Promise.all(userRequests).then((responses) => {
      let newUsers: User[] = []
      for (const response of responses) {
        if (response.Success && response.Result) {
          newUsers.push(response.Result)
        }
      }
      setUsers(newUsers)
    })
  }, [usernames])
  return (
    <>
      {users.length > 0 && (
        <div style={{ display: 'flex', justifyContent: 'center' }}>
          <Typography variant='h3' color='text.primary'>
            Users
          </Typography>
        </div>
      )}
      {users.map((user) => (
        <Collapse key={user.Username} in appear>
          <Card key={user.Username} style={{ margin: '10px' }}>
            <CardActionArea
              sx={{
                flex: '1 0 auto',
              }}
            >
              <CardContent
                sx={{
                  display: 'flex',
                  alignItems: 'center',
                  alignContent: 'center',
                  marginBottom: '5px',
                  justifyContent: 'space-between',
                }}
                component='div'
              >
                <Box
                  sx={{
                    display: 'flex',
                    alignItems: 'center',
                    alignContent: 'center',
                    marginBottom: '5px',
                  }}
                >
                  <Avatar style={{ marginRight: '10px' }}>
                    {initials(user.Username)}
                  </Avatar>
                  <Typography
                    variant='subtitle1'
                    color='text.secondary'
                    component='div'
                  >
                    {user.Username}
                  </Typography>
                </Box>
                <Box
                  sx={{
                    display: 'flex',
                    alignItems: 'center',
                    alignContent: 'center',
                    marginBottom: '5px',
                    gap: '5px',
                  }}
                >
                  <Typography variant='body1' color='text.secondary'>
                    Posts: {user.Posts ? user.Posts.length : 0},
                  </Typography>
                  <Typography variant='body1' color='text.secondary'>
                    Comments: {user.Comments ? user.Comments.length : 0},
                  </Typography>
                  <Typography variant='body1' color='text.secondary'>
                    Voted Comments:{' '}
                    {(user.UPVotedComments ? user.UPVotedComments.length : 0) +
                      (user.DOWNVotedComments
                        ? user.DOWNVotedComments.length
                        : 0)}
                    ,
                  </Typography>
                  <Typography variant='body1' color='text.secondary'>
                    Voted Posts:{' '}
                    {(user.UPVotedPosts ? user.UPVotedPosts.length : 0) +
                      (user.DOWNVotedPosts ? user.DOWNVotedPosts.length : 0)}
                  </Typography>
                </Box>
              </CardContent>
            </CardActionArea>
          </Card>
        </Collapse>
      ))}
    </>
  )
}

function SearchResultPosts({
  currentUserState,
  postIds,
}: {
  currentUserState: CurrentUserState
  postIds: number[]
}) {
  const [posts, setPosts] = useState<Post[]>([])
  useEffect(() => {
    getPosts(postIds).then((res) => {
      if (res.Success && res.Result) {
        let resArray: Post[] = res.Result
        if (resArray.length > 5) {
          setPosts(resArray.slice(0, 5))
        } else {
          setPosts(res.Result)
        }
      } else {
        setPosts([])
      }
    })
  }, [postIds])
  return (
    <>
      {posts.length > 0 && (
        <div style={{ display: 'flex', justifyContent: 'center' }}>
          <Typography variant='h3' color='text.primary'>
            Posts
          </Typography>
        </div>
      )}
      {posts.map((post) => (
        <Collapse
          key={post.Url + post.Creator + post.Date}
          in={true}
          appear={true}
        >
          <Card
            sx={{ display: 'flex', justifyContent: 'space-between' }}
            style={{ margin: '10px' }}
          >
            <Box sx={{ display: 'flex', flexDirection: 'column' }}>
              <CardContent sx={{ flex: '1 0 auto' }} component='div'>
                <Box
                  sx={{
                    display: 'flex',
                    alignItems: 'center',
                    alignContent: 'center',
                    marginBottom: '5px',
                  }}
                >
                  <Avatar style={{ marginRight: '10px' }}>
                    {initials(post.Creator)}
                  </Avatar>
                  <Typography
                    variant='subtitle1'
                    color='text.secondary'
                    component='div'
                  >
                    {post.Creator}
                  </Typography>
                </Box>
                <Typography component='div' variant='h5'>
                  {post.Caption}
                </Typography>
              </CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', pl: 1, pb: 1 }}>
                <IconButton
                  onClick={() =>
                    handleVotePost(
                      post,
                      null,
                      currentUserState,
                      isPostUpVoted(post.Id, currentUserState) ? 0 : 1
                    )
                  }
                >
                  <KeyboardArrowUp
                    color={
                      isPostUpVoted(post.Id, currentUserState)
                        ? 'secondary'
                        : 'action'
                    }
                  />
                </IconButton>
                <IconButton
                  onClick={() =>
                    handleVotePost(
                      post,
                      null,
                      currentUserState,
                      isPostDownVoted(post.Id, currentUserState) ? 0 : -1
                    )
                  }
                >
                  <KeyboardArrowDown
                    color={
                      isPostDownVoted(post.Id, currentUserState)
                        ? 'secondary'
                        : 'action'
                    }
                  />
                </IconButton>
              </Box>
            </Box>
            {filetype(post.Url) == 'mp4' || filetype(post.Url) == 'odd' ? (
              <CardMedia
                controls={false}
                src={post.Url}
                style={{ maxWidth: 150, maxHeight: 160 }}
                component='video'
              />
            ) : (
              <CardMedia
                src={post.Url}
                component='img'
                style={{ maxWidth: 150, maxHeight: 160 }}
              />
            )}
          </Card>
        </Collapse>
      ))}
    </>
  )
}

export default Search
