import { Button, Card, CardActions, CardContent, CardMedia, CircularProgress, Grid, Typography, } from "@mui/material";
import { useState } from "react";
import { PostProps } from "../utils/types"
import { savePost } from "../utils/serverFunctions";

function Post({ loading, caption, src, setRefresh, random: controls, comments }: PostProps) {
  const [saving, setSaving] = useState(false)
  const [saveMenu, setSaveMenu] = useState<boolean>(false)

  return <Grid container justifyContent="center" marginBottom={4}>
    <Grid item>
      <Card>
        {loading ?
          <Grid
            container
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
          <Button
            variant="contained"
            style={{ backgroundColor: "#EC407A", color: "white" }}
            fullWidth={true}
            onClick={() => {
              if (setRefresh) setRefresh(true)
            }}>
            Pass
          </Button>
          <Button
            variant="contained"
            style={{ backgroundColor: "#66BB6A", color: "white" }}
            fullWidth={true}
            onClick={() => {
              setSaveMenu(true)
              setSaving(true)
              savePost(src, "").then((res) => {
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
        {comments &&
          <Button
            variant="contained"
            style={{
              margin: "8px",
              width: "calc(100% - 16px)"
            }}>
            Comments
          </Button>
        }
      </Card>
    </Grid>
  </Grid>
}

function filetype(src: string) {
  if (src) return src.substring(src.lastIndexOf(".") + 1)
}

export default Post
