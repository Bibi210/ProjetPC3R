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
import { Post as PostType, User } from "../utils/types"
import { getPosts, getPrivateProfile } from "../utils/serverFunctions";

function Profile() {
  const [loading, setLoading] = useState(true)
  const [user, setUser] = useState<User | null>(null)
  const [posts, setPosts] = useState<PostType[]>([])

  useEffect(() => {
    getPrivateProfile().then(userRes => {
      if (userRes.Success) {
        let newUser: User = userRes.Result
        setUser(newUser)
        getPosts(newUser.Posts).then(postsRes => {
          if (postsRes.Success) {
            console.log(postsRes.Result)
            if (postsRes.Result == null) {
              setPosts([])
            } else {
              setPosts(postsRes.Result)
            }
          } else {
            console.log(postsRes.Message)
            // notify
          }
        })
      }
      setLoading(false)
    })
  }, [])
  return <Container>
    <Card>
      <CardHeader
        avatar={
          loading ? <CircularProgress /> :
            user ? <Avatar> {initials(user.Username)}</Avatar> : <Avatar></Avatar>
        }
        title={loading ? "" : user ? <Typography variant="h6">{user.Username}</Typography> : "No connected user"}
        action={loading ? <CircularProgress /> : user &&
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
          user ?
            <List>
              {posts.length == 0 &&
                <Grid container justifyContent="center">
                  <Grid item>
                    <Typography variant="h5">Nothing here yet</Typography>
                  </Grid>
                </Grid>
              }
              {posts.map((post: PostType) =>
                <ListItem key={post.Url + post.Creator + post.Date}>
                  <Post loading={false} src={post.Url} caption={post.Caption} random={false} comments={true} />
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