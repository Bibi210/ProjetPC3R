import { useEffect, useState } from "react";
import { getComments, saveComment } from "../utils/serverFunctions";
import { Comment, getCurrentUser } from "../utils/types"
import {
  Avatar,
  Box,
  Button,
  Card,
  CardActionArea,
  CardContent,
  CircularProgress,
  Collapse,
  Container,
  Grid,
  List,
  TextField,
  Typography
} from "@mui/material";
import { ArrowUpward } from "@mui/icons-material";
import { initials } from "./Profile";

function Comments({ postId, ids, showComments }: { postId: number, ids: number[], showComments: boolean }) {
  const [comments, setComments] = useState<Comment[]>([]);
  const [show, setShow] = useState(false)
  const [newComment, setNewComment] = useState("");
  const [savingNewComment, setSavingNewComment] = useState(false);
  let currentUser = getCurrentUser()

  function handleSaveComment() {
    setSavingNewComment(true)
    saveComment(postId, newComment).then(res => {
      setSavingNewComment(false)
      if (res.Success) {
        setComments([...comments, {
          Msg: {
            Content: newComment,
            Date: new Date().toString(),
            Sender: currentUser ? currentUser.Username : "should not happen"
          },
          Upvotes: 0
        }])
        setNewComment("")
      } else {
        // !notify
      }
    })
  }

  useEffect(() => {
    if (showComments) {
      setShow(false)
      getComments(ids).then(res => {
        if (res.Success && res.Result) {
          setComments(res.Result)
          setShow(true)
        }
      })
    }
  }, [ids, showComments])


  return <Container>
    {comments.length == 0 ?
      <Collapse in={showComments} appear={showComments}>
        <Grid container justifyContent="center">
          <Grid item>
            <Typography variant="h5">Nothing here yet</Typography>
          </Grid>
        </Grid>
      </Collapse>
      :
      <List sx={{ maxHeight: "500px", overflow: "scroll" }}>
        {comments.map((comment) =>
          <Collapse key={comment.Msg.Sender + comment.Msg.Date} appear={showComments && show} in={showComments && show}>
            <Card variant="outlined" sx={{ marginTop: "10px" }}>
              <CardActionArea>
                <Box sx={{ display: 'flex' }}>
                  <CardContent sx={{ flex: '1 0 auto' }} component="div">
                    <Box sx={{ display: 'flex', alignItems: "center", alignContent: "center", marginBottom: "15px" }}>
                      <Avatar style={{ marginRight: "10px" }}>{initials(comment.Msg.Sender)}</Avatar>
                      <Typography variant="subtitle1" color="text.secondary" component="div">
                        {comment.Msg.Sender}
                      </Typography>
                    </Box>
                    <Typography component="div" variant="body1">
                      {comment.Msg.Content}
                    </Typography>
                  </CardContent>
                  <Box sx={{ display: 'flex', alignItems: 'center', marginRight: "25px" }}>
                    <ArrowUpward />
                  </Box>
                </Box>
              </CardActionArea>
            </Card>
          </Collapse>
        )}
      </List>
    }
    {currentUser &&
      <Collapse
        in={showComments}
        appear={showComments}
      >
        <Box sx={{ display: 'flex', alignContent: "center", alignItems: "center", gap: "15px", margin: "15px 0" }}>
          <TextField
            variant="outlined"
            value={newComment}
            onKeyUp={(e) => {
              if (e.key == "Enter") {
                handleSaveComment()
              }
            }}
            onChange={(e) => {
              setNewComment(e.target.value)
            }}
            label="Comment content"
            fullWidth
          />
          <Button
            variant="contained"
            size="small"
            onClick={() => {
              handleSaveComment()
            }}
          >
            {savingNewComment ? <CircularProgress /> : "Add a new comment"}
          </Button>
        </Box>
      </Collapse>
    }
  </Container>
}

export default Comments