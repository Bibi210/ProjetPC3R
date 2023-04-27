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


type Post = {
    Caption: string,
    Creator: string,
    Date: string,
    Upvotes: number,
    Url: string
}

export async function getPrivateProfile() {
    let res = await fetch(window.location.origin + "/api/get_private_profile")
    let json = await res.json()
    console.log(json)
    return json
}

async function getSavedShitpost(id:number) {
    let res = await fetch(window.location.origin + "/api/get_saved_shitpost", {

    })
}

function initials(name: string) {
    return name.split(" ").map((n) => n[0].toUpperCase()).join()
}

function Profile() {
    const [loading, setLoading] = useState(true)
    const [user, setUser] = useState<any>(null)

    useEffect(() => {
        getPrivateProfile().then(res => {
            setUser(res.Result)
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
                title={user ? user.Username : "No connected user"}
                action={loading ? <CircularProgress /> :
                    user &&
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
                            {user.Posts && user.Posts.map((post: Post) =>
                                <ListItem>
                                    <Post loading={false} src={post.Url} caption={post.Caption} controls={false} />
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

export default Profile
