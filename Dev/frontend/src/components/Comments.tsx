import { useEffect, useState } from 'react'
import { getComments, saveComment, voteComment } from '../utils/serverFunctions'
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
import { KeyboardArrowDown, KeyboardArrowUp } from '@mui/icons-material'
import { initials } from './Profile'
import {
  handleVoteComment,
  isCommentDownVoted,
  isCommentUpVoted,
} from '../utils/utils'

function Comments({
  currentUserState,
  post,
  showComments,
}: CommentComponentProps) {
  const [comments, setComments] = useState<Comment[]>([])
  const [show, setShow] = useState(false)
  const [newComment, setNewComment] = useState('')
  const [savingNewComment, setSavingNewComment] = useState(false)

  function handleSaveComment() {
    setSavingNewComment(true)
    saveComment(post.Id, newComment).then((res) => {
      setSavingNewComment(false)
      if (res.Success && res.Result) {
        if (!currentUserState || !currentUserState.get) {
          alert('could not retrieve current user')
        } else {
          let newDate = new Date()
          setComments([
            ...comments,
            {
              Id: res.Result.Id,
              Msg: {
                Content: newComment,
                Date:
                  newDate.toString().split(' ').slice(0, 2).join(' ') +
                  ' ' +
                  newDate.getDate() +
                  ' ' +
                  newDate.getHours().toString().padStart(2, '0') +
                  ':' +
                  (newDate.getMinutes() % 60).toString().padStart(2, '0') +
                  ':' +
                  (newDate.getSeconds() % 3600).toString().padStart(2, '0') +
                  ' ' +
                  newDate.getFullYear().toString().padStart(4, '0'),
                Sender: currentUserState.get.Username,
              },
              Upvotes: 0,
            },
          ])
          setShow(true)
        }
        setNewComment('')
      } else {
        alert(res.Message)
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
                      comments,
                      setComments,
                      currentUserState
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
                        <Typography
                          variant='subtitle2'
                          color='text.secondary'
                          style={{ marginLeft: '10px' }}
                        >
                          {comment.Msg.Date}
                        </Typography>
                      </Box>
                      <Typography component='div' variant='body1'>
                        {comment.Msg.Content}
                      </Typography>
                    </CardContent>
                    <Box
                      sx={{
                        display: 'flex',
                        alignItems: 'center',
                        alignContent: 'center',
                        justifyContent: 'center',
                        marginRight: '25px',
                      }}
                    >
                      <KeyboardArrowUp
                        color={
                          isCommentUpVoted(comment.Id, currentUserState)
                            ? 'error'
                            : 'action'
                        }
                      />
                      <Typography
                        component='div'
                        variant='body1'
                        style={{ marginTop: '4px' }}
                      >
                        {comment.Upvotes}
                      </Typography>
                      <KeyboardArrowDown
                        color={
                          isCommentDownVoted(comment.Id, currentUserState)
                            ? 'error'
                            : 'action'
                        }
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
