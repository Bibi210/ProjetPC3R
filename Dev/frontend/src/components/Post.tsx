import {
    Button,
    Card,
    CardActions,
    CardMedia,
    CircularProgress,
    Grid,
} from "@mui/material";
import React, {Dispatch, SetStateAction} from "react";

type PostProps = {
    loading: boolean,
    src: string,
    setRefresh: Dispatch<SetStateAction<boolean>>
}

function Post({loading, src, setRefresh}: PostProps) {
    return <Grid container justifyContent="center">
        <Grid item>
            <Card>
                {loading ?
                    <Grid container
                          justifyContent="center"
                          alignContent="center"
                          style={{width: "800px", height: "600px"}}>
                        <Grid item><CircularProgress/></Grid>
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
                <CardActions>
                    <Button variant="contained"
                            style={{backgroundColor: "#EC407A", color: "white"}}
                            fullWidth={true}
                            onClick={() => setRefresh(true)}>
                        Pass
                    </Button>
                    <Button variant="contained"
                            style={{backgroundColor: "#66BB6A", color: "white"}}
                            fullWidth={true}
                            onClick={() => setRefresh(true)}>
                        Save
                    </Button>
                </CardActions>
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

function filetype(src: string) {
    if (src)
        return src.substring(src.lastIndexOf(".") + 1)
}

export default Post