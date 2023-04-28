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
      let newUser: User = userRes.Result
      setUser(newUser)
      getPosts(newUser.Posts).then(postsRes => {
        if (postsRes.Success) {
          setPosts(postsRes.Result)
        } else {
          console.log(postsRes.Message)
          // notifie
        }
      })
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
        title={loading ? "" : user ? user.Username : "No connected user"}
        action={loading ? <CircularProgress /> : user &&
          <Grid container spacing={1}>
            <Grid item>
              <Button variant="contained"
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
              {posts.map((post: PostType) =>
                <ListItem>
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

function initials(name: string): string {
  return name.split(" ").map((n) => n[0].toUpperCase()).join()
}

export default Profile