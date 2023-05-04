import {
  Avatar,
  Button,
  Card,
  CardContent,
  CardHeader,
  CircularProgress,
  Container,
  Grid,
  List,
  ListItem,
  Typography,
} from "@mui/material";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import Post from "./Post";
import { CurrentUserState, Post as PostType } from "../utils/types"
import { getPosts } from "../utils/serverFunctions";

function Profile({ currentUserState }: { currentUserState: CurrentUserState }) {
  const [loading, setLoading] = useState(false)
  const [posts, setPosts] = useState<PostType[]>([])

  useEffect(() => {
    console.log("currentUserState", currentUserState)
    if (currentUserState && currentUserState.get) {
      getPosts(currentUserState.get.Posts).then(postsRes => {
        if (postsRes.Success) {
          setPosts(postsRes.Result ? postsRes.Result : [])
        } else {
          console.log(postsRes.Message)
          // notify
        }
        setLoading(false)
      })
    }
  }, [])
  return <Container>
    <Card>
      <CardHeader
        avatar={
          loading ? <CircularProgress /> :
            currentUserState && currentUserState.get ? <Avatar> {initials(currentUserState.get.Username)}</Avatar> :
              <Avatar></Avatar>
        }
        title={loading ? "" : currentUserState && currentUserState.get ?
          <Typography variant="h6">{currentUserState.get.Username}</Typography> : "No connected user"}
        action={loading ? <CircularProgress /> : currentUserState && currentUserState.get &&
          <Grid container spacing={1}>
            <Grid item>
              <Button
                variant="contained"
                style={{ backgroundColor: "#EF5350", color: "white" }}>
                Delete account
              </Button>
            </Grid>
            <Grid item>
              <Link to="/logout">
                <Button variant="contained">Logout</Button>
              </Link>
            </Grid>
          </Grid>
        }
      ></CardHeader>

      <CardContent>
        {loading ? <CircularProgress /> :
          currentUserState && currentUserState.get ?
            <List>
              {posts.length == 0 ?
                <Grid container justifyContent="center">
                  <Grid item>
                    <Typography variant="h5">Nothing here yet</Typography>
                  </Grid>
                </Grid>
                :
                posts.map((post: PostType) =>
                  <ListItem key={post.Url + post.Creator + post.Date}>
                    <Post
                      currentUserState={currentUserState}
                      loading={false}
                      post={post}
                      randomMode={false}
                    />
                  </ListItem>
                )}
            </List>
            :
            <Link to="/login">
              <Button fullWidth style={{ backgroundColor: "#3F51B5", color: "white" }}>
                Login
              </Button>
            </Link>
        }
      </CardContent>
    </Card>
  </Container>
}

export function initials(name: string): string {
  return name.split(" ").map((n) => n[0].toUpperCase()).join()
}

export default Profile