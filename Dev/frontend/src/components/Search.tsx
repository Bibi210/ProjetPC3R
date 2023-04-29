import {
  Avatar,
  Box,
  Card,
  CardActionArea,
  CardContent,
  CardMedia,
  Container,
  IconButton,
  InputAdornment,
  TextField,
  Typography
} from "@mui/material";
import { useEffect, useRef, useState } from "react";
import { getPosts, getPublicProfile, search } from "../utils/serverFunctions";
import { Post, SearchResults, ServerResponse, User } from "../utils/types";
import { initials } from "./Profile";
import { filetype } from "./Post";
import { ArrowUpward, Chat } from "@mui/icons-material";


function Search() {
  const [searchResults, setSearchResults] = useState<SearchResults>({ ShitPosts: [], Users: [] });
  let keyTimeout = useRef(0)

  function onChange(value: string) {
    if (keyTimeout.current != 0) {
      clearTimeout(keyTimeout.current)
    }
    keyTimeout.current = setTimeout(() => {
      if (value != "") {
        search(value).then(searchRes => {
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
      }
    }, 500)
  }

  return <Container>
    <TextField
      fullWidth={true}
      onChange={(e) => onChange(e.target.value)}
      InputProps={{
        startAdornment: (
          <InputAdornment position="start">
            🔍 search
          </InputAdornment>
        ),
      }}
      style={{ marginBottom: "20px" }}
    />
    <SearchResultUsers usernames={searchResults.Users} />
    <SearchResultPosts postIds={searchResults.ShitPosts} />
  </Container>
}

function SearchResultUsers({ usernames }: { usernames: string[] }) {
  const [users, setUsers] = useState<User[]>([]);
  useEffect(() => {
    let userRequests: Promise<ServerResponse>[] = []
    for (const username of usernames) {
      userRequests.push(getPublicProfile(username))
    }
    Promise.all(userRequests).then(responses => {
      let newUsers: User[] = []
      for (const response of responses) {
        if (response.Success && response.Result) {
          newUsers.push(response.Result)
        }
      }
      setUsers(newUsers)
    })
  }, [usernames])
  return <>
    {users.length > 0 && <Typography variant="h3">Users</Typography>}
    {users.map(user =>
      <Card key={user.Username} style={{ margin: "10px" }}>
        <CardActionArea sx={{
          flex: '1 0 auto'
        }}>
          <CardContent sx={{
            display: 'flex',
            alignItems: "center",
            alignContent: "center",
            marginBottom: "5px",
            justifyContent: "space-between"
          }} component="div">
            <Box sx={{ display: 'flex', alignItems: "center", alignContent: "center", marginBottom: "5px" }}>
              <Avatar style={{ marginRight: "10px" }}>{initials(user.Username)}</Avatar>
              <Typography variant="subtitle1" color="text.secondary" component="div">
                {user.Username}
              </Typography>
            </Box>
            <Box sx={{
              display: 'flex',
              alignItems: "center",
              alignContent: "center",
              marginBottom: "5px",
              gap: "5px"
            }}>
              <Typography variant="body1" color="text.secondary">
                Posts: {user.Posts ? user.Posts.length : "0"},
              </Typography>
              <Typography variant="body1" color="text.secondary">
                Comments: {user.Comments ? user.Comments.length : "0"},
              </Typography>
              <Typography variant="body1" color="text.secondary">
                Voted Comments: {user.VotedComments ? user.VotedComments.length : "0"},
              </Typography>
              <Typography variant="body1" color="text.secondary">Voted
                Posts: {user.VotedPosts ? user.VotedPosts.length : "0"}
              </Typography>
            </Box>
          </CardContent>
        </CardActionArea>
      </Card>
    )}
  </>
}

function SearchResultPosts({ postIds }: { postIds: number[] }) {
  const [posts, setPosts] = useState<Post[]>([]);
  useEffect(() => {
    getPosts(postIds).then(res => {
      if (res.Success && res.Result) {
        let resArray: Post[] = res.Result
        if (resArray.length > 5) {
          setPosts(resArray.slice(0, 5))
        } else {
          setPosts(res.Result)
        }
      }
    })
  }, [postIds])
  return <>
    {posts.length > 0 && <Typography variant="h3">Posts</Typography>}
    {posts.map(post =>
      <Card
        key={post.Url + post.Creator + post.Date}
        sx={{ display: 'flex', justifyContent: 'space-between' }}
        style={{ margin: "10px" }}>
        <Box sx={{ display: 'flex', flexDirection: 'column' }}>
          <CardContent sx={{ flex: '1 0 auto' }} component="div">
            <Box sx={{ display: 'flex', alignItems: "center", alignContent: "center", marginBottom: "5px" }}>
              <Avatar style={{ marginRight: "10px" }}>{initials(post.Creator)}</Avatar>
              <Typography variant="subtitle1" color="text.secondary" component="div">
                {post.Creator}
              </Typography>
            </Box>
            <Typography component="div" variant="h5">
              {post.Caption}
            </Typography>
          </CardContent>
          <Box sx={{ display: 'flex', alignItems: 'center', pl: 1, pb: 1 }}>
            <IconButton aria-label="upvote">
              <ArrowUpward />
            </IconButton>
            <IconButton aria-label="comments">
              <Chat />
            </IconButton>
          </Box>
        </Box>
        {
          (filetype(post.Url) == "mp4" || filetype(post.Url) == "odd") ?
            <CardMedia
              controls={false}
              src={post.Url}
              style={{ maxWidth: 150, maxHeight: 160 }}
              component="video"
            />
            :
            <CardMedia
              src={post.Url}
              component="img"
              style={{ maxWidth: 150, maxHeight: 160 }}
            />
        }
      </Card>
    )}
  </>
}

export default Search
