import { useEffect, useState } from 'react'
import {
  getComments,
  saveComment,
  upvoteComment,
} from '../utils/serverFunctions'
import { Comment, CommentComponentProps } from '../utils/types'
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
  Typography,
} from '@mui/material'
import { ArrowUpward } from '@mui/icons-material'
import { initials } from './Profile'

function Comments({
  currentUserState,
  post,
  showComments,
}: CommentComponentProps) {
  const [comments, setComments] = useState<Comment[]>([])
  const [show, setShow] = useState(false)
  const [newComment, setNewComment] = useState('')
  const [savingNewComment, setSavingNewComment] = useState(false)

  function isVoted(commentId: number) {
    if (!currentUserState || !currentUserState.get) {
      // !notify
      return false
    } else if (currentUserState.get.VotedComments) {
      return currentUserState.get.VotedComments.includes(commentId)
    } else {
      return false
    }
  }

  function handleVoteComment(comment: Comment, up: boolean) {
    upvoteComment(comment.Id, up ? 1 : -1).then((res) => {
      if (res.Success) {
        let newComments: Comment[] = []
        for (const c of comments) {
          if (c.Id == comment.Id) {
            newComments.push({
              Id: c.Id,
              Msg: c.Msg,
              Upvotes: c.Upvotes + (up ? 1 : -1),
            })

            if (currentUserState.get) {
              let newVotedComments
              if (up)
                newVotedComments = currentUserState.get.VotedComments
                  ? [...currentUserState.get.VotedComments, comment.Id]
                  : [comment.Id]
              else
                newVotedComments = currentUserState.get.VotedComments
                  ? currentUserState.get.VotedComments.filter(
                      (comId) => comId != c.Id
                    )
                  : []

              currentUserState.set({
                ...currentUserState.get,
                VotedComments: newVotedComments,
              })
            } else {
              console.error(
                'upvoted comment without active user, should not happen'
              )
            }
          } else {
            newComments.push(c)
          }
        }
        console.log('arrived to setcomment')

        setComments(newComments)
      }
    })
  }

  function handleSaveComment() {
    setSavingNewComment(true)
    saveComment(post.Id, newComment).then((res) => {
      setSavingNewComment(false)
      if (res.Success && res.Result) {
        if (!currentUserState || !currentUserState.get) {
          // !notify
        } else {
          setComments([
            ...comments,
            {
              Id: res.Result.Id,
              Msg: {
                Content: newComment,
                Date: new Date().toString(),
                Sender: currentUserState.get.Username,
              },
              Upvotes: 0,
            },
          ])
          setShow(true)
        }
        setNewComment('')
      } else {
        // !notify
      }
    })
  }

  useEffect(() => {
    if (showComments) {
      setShow(false)
      getComments(post.CommentIds).then((res) => {
        if (res.Success && res.Result) {
          setComments(res.Result)
          setShow(true)
        }
      })
    }
  }, [post, showComments])

  return (
    <Container>
      {comments.length == 0 ? (
        <Collapse in={showComments} appear={showComments}>
          <Grid container justifyContent='center'>
            <Grid item>
              <Typography variant='h5'>Nothing here yet</Typography>
            </Grid>
          </Grid>
        </Collapse>
      ) : (
        <List
          sx={{ maxHeight: '500px', overflowY: 'auto', overflowX: 'hidden' }}
        >
          {comments.map((comment) => (
            <Collapse
              key={comment.Msg.Sender + comment.Msg.Date}
              appear={showComments && show}
              in={showComments && show}
            >
              <Card variant='outlined' sx={{ marginTop: '10px' }}>
                <CardActionArea
                  onClick={() =>
                    handleVoteComment(
                      comment,
                      !isVoted(comment.Id)
                    )
                  }
                >
                  <Box sx={{ display: 'flex' }}>
                    <CardContent sx={{ flex: '1 0 auto' }} component='div'>
                      <Box
                        sx={{
                          display: 'flex',
                          alignItems: 'center',
                          alignContent: 'center',
                          marginBottom: '15px',
                        }}
                      >
                        <Avatar style={{ marginRight: '10px' }}>
                          {initials(comment.Msg.Sender)}
                        </Avatar>
                        <Typography
                          variant='subtitle1'
                          color='text.secondary'
                          component='div'
                        >
                          {comment.Msg.Sender}
                        </Typography>
                      </Box>
                      <Typography component='div' variant='body1'>
                        {comment.Id + ' ' + comment.Msg.Content}
                      </Typography>
                    </CardContent>
                    <Box
                      sx={{
                        display: 'flex',
                        alignItems: 'center',
                        marginRight: '25px',
                      }}
                    >
                      <ArrowUpward
                        color={isVoted(comment.Id) ? 'primary' : 'success'}
                      />
                    </Box>
                  </Box>
                </CardActionArea>
              </Card>
            </Collapse>
          ))}
        </List>
      )}
      {currentUserState && currentUserState.get && (
        <Collapse in={showComments} appear={showComments}>
          <Box
            sx={{
              display: 'flex',
              alignContent: 'center',
              alignItems: 'center',
              gap: '15px',
              margin: '15px 0',
            }}
          >
            <TextField
              variant='outlined'
              value={newComment}
              onKeyUp={(e) => {
                if (e.key == 'Enter') {
                  handleSaveComment()
                }
              }}
              onChange={(e) => {
                setNewComment(e.target.value)
              }}
              label='Comment content'
              fullWidth
            />
            <Button
              variant='contained'
              size='small'
              onClick={() => {
                handleSaveComment()
              }}
            >
              {savingNewComment ? <CircularProgress /> : 'Add a new comment'}
            </Button>
          </Box>
        </Collapse>
      )}
    </Container>
  )
}

export default Comments
