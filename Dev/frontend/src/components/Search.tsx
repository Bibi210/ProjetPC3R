import {
  Avatar,
  Box,
  Card,
  CardContent,
  CardMedia,
  Container,
  IconButton,
  InputAdornment,
  TextField,
  Typography
} from "@mui/material";
import React, { useEffect, useRef, useState } from "react";
import { getPosts, search } from "../utils/serverFunctions";
import { Post, SearchResults, User } from "../utils/types";
import { initials } from "./Profile";
import { filetype } from "./Post";


function Search() {

  const [searchResults, setSearchResults] = useState<SearchResults>({ ShitPosts: [], Users: [] });
  const [inputValue, setInputValue] = useState("");

  let keyTimeout = useRef(0)

  function onChange(value: string) {
    setInputValue(value)
    if (keyTimeout.current != 0) {
      clearTimeout(keyTimeout.current)
    }
    keyTimeout.current = setTimeout(() => {
      search(inputValue).then(searchRes => {
        console.log(searchRes)
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
    }, 500)
  }

  return <Container>
    <TextField
      fullWidth={true}
      value={inputValue}
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
  console.log("usernames", usernames)
  const [users, setUsers] = useState<User[]>([]);
  useEffect(() => {
  }, [usernames])
  return <>
    {users.length > 0 && <Typography variant="h3">Users</Typography>}
  </>
}

function SearchResultPosts({ postIds }: { postIds: number[] }) {
  const [posts, setPosts] = useState<Post[]>([]);
  useEffect(() => {
    getPosts(postIds).then(res => {
      if (res.Success && res.Result) {
        console.log("posts", res.Result)
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
      <Card sx={{ display: 'flex', justifyContent: 'space-between' }}>
        <Box sx={{ display: 'flex', flexDirection: 'column' }}>
          <CardContent sx={{ flex: '1 0 auto' }}>
            <Box sx={{ display: 'flex', alignItems: "center", alignContent: "center", marginBottom: "5px" }}>
              <Avatar style={{marginRight: "10px"}}>{initials(post.Creator)}</Avatar>
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
              <span className="material-icons">arrow_upward</span>
            </IconButton>
            <IconButton aria-label="comments">
              <span className="material-icons">chat</span>
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