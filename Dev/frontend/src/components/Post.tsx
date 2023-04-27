import {
    Button,
    Card,
    CardActions,
    CardContent,
    CardMedia,
    CircularProgress,
    Grid,
    Typography,
} from "@mui/material";
import { Dispatch, SetStateAction, useState } from "react";

type PostProps = {
    loading: boolean,
    caption: string,
    src: string,
    setRefresh?: Dispatch<SetStateAction<boolean>>,
    controls: boolean
}

function Post({ loading, caption, src, setRefresh, controls }: PostProps) {
    const [saving, setSaving] = useState(false)
    return <Grid container justifyContent="center">
        <Grid item>
            <Card>
                {loading ?
                    <Grid container
                        justifyContent="center"
                        alignContent="center"
                        style={{ width: "800px", height: "600px" }}>
                        <Grid item><CircularProgress /></Grid>
                    </Grid>
                    :
                    (filetype(src) == "mp4" || filetype(src) == "odd") ?
                        <CardMedia
                            controls={true}
                            src={src}
                            component="video"
                            style={{
                                width: "800px",
                                height: "600px"
                            }}
                        />
                        :
                        <CardMedia
                            src={src}
                            component="img"
                            style={{
                                width: "800px",
                                height: "600px"
                            }}
                        />
                }
                {caption != "" &&
                    <CardContent>
                        <Typography variant="body2">{caption}</Typography>
                    </CardContent>}
                {controls && <CardActions>
                    <Button variant="contained"
                        style={{ backgroundColor: "#EC407A", color: "white" }}
                        fullWidth={true}
                        onClick={() => {
                            if (setRefresh) setRefresh(true)
                        }}>
                        Pass
                    </Button>
                    <Button variant="contained"
                        style={{ backgroundColor: "#66BB6A", color: "white" }}
                        fullWidth={true}
                        onClick={() => {
                            setSaving(true)
                            savePost(src).then((res) => {
                                setSaving(false)
                                if (setRefresh) {
                                    res.Success ?
                                        setRefresh(true) :
                                        console.log(res.Message)
                                }
                            })
                        }}>
                        {saving ? <CircularProgress /> : "Save"}
                    </Button>
                </CardActions>}
                <Button
                    variant="contained"
                    style={{
                        margin: "8px",
                        width: "calc(100% - 16px)"
                    }}>
                    Comments
                </Button>
            </Card>
        </Grid>
    </Grid>
}

async function savePost(url: string) {
    url = url.replaceAll(" ", "_")
    let caption = url.substring(url.lastIndexOf("/") + 1)
    let req = await fetch(window.location.origin + "/api/save_shitpost", {
        method: "POST",
        body: `{"url":"${url}", "caption": "${caption}"}`
    })
    return await req.json()
}

function filetype(src: string) {
    if (src) return src.substring(src.lastIndexOf(".") + 1)
}

export default Post
