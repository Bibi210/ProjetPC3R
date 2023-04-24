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
} from "@mui/material";
import React, {useEffect, useState} from "react";
import {Link} from "react-router-dom";


export async function getPrivateProfile() {
    let res = await fetch(window.location.origin + "/api/get_private_profile")
    return await res.json()
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
                    loading ? <CircularProgress/> :
                        user ? <Avatar> {initials(user.Username)}</Avatar> : <Avatar></Avatar>
                }
                title={user ? user.Username : ""}
                action={loading ? <CircularProgress/> :
                    user ?
                        <Grid container spacing={1}>
                            <Grid item>
                                <Button variant="contained"
                                        style={{backgroundColor: "#EF5350", color: "white"}}>
                                    Delete account
                                </Button>
                            </Grid>
                            <Grid item>
                                <Link to="/logout">
                                    <Button variant="contained">Logout</Button>
                                </Link>
                            </Grid>
                        </Grid>
                        :
                        <Link to="/login">
                            <Button style={{backgroundColor: "#3F51B5", color: "white"}}>
                                Login
                            </Button>
                        </Link>
                }
            ></CardHeader>

            <CardContent>
                {loading ? <CircularProgress/> :
                    <List>
                        {/*<ListItem><Post title="Title Post 1" contentText="post 1 content" hasSrc={false}/></ListItem>*/}
                        {/*<ListItem><Post title="Title Post 2" contentText="post 2 content" hasSrc={false}/></ListItem>*/}
                    </List>
                }
            </CardContent>
        </Card>
    </Container>
}

export default Profile
