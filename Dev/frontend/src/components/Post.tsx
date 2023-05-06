import {
  Box,
  Button,
  Card,
  CardActions,
  CardContent,
  CardMedia,
  CircularProgress,
  Grid,
  IconButton,
  Menu,
  TextField,
  Typography,
} from '@mui/material'
import React, { useState } from 'react'
import { PostComponentProps } from '../utils/types'
import { ArrowUpward, Chat, Check } from '@mui/icons-material'
import { savePost } from '../utils/serverFunctions'
import Comments from './Comments'
import { handleUpvotePost, isVoted } from '../utils/utils'

function Post({
  currentUserState,
  loading,
  post,
  setRefresh,
  randomMode,
}: PostComponentProps) {
  const [saving, setSaving] = useState(false)
  const [saveMenuAnchor, setSaveMenuAnchor] = useState<null | HTMLElement>(null)
  const [savingCaption, setSavingCaption] = useState('')
  const [showComments, setShowComments] = useState(false)
  let openSaveMenu = Boolean(saveMenuAnchor)

  function handleSaveBtnClick(event: React.MouseEvent<HTMLElement>) {
    setSaveMenuAnchor(event.currentTarget)
  }

  function handleSavePost() {
    savePost(post.Url, savingCaption).then((res) => {
      setSavingCaption('')
      setSaveMenuAnchor(null)
      setSaving(false)
      if (setRefresh) {
        res.Success ? setRefresh(true) : alert(res.Message)
      }
    })
  }

  return (
    <Grid container justifyContent='center' marginBottom={4}>
      <Grid item>
        <Card>
          {loading ? (
            <Grid
              container
              justifyContent='center'
              alignContent='center'
              style={{ width: '800px', height: '600px' }}
            >
              <Grid item>
                <CircularProgress />
              </Grid>
            </Grid>
          ) : filetype(post.Url) == 'mp4' || filetype(post.Url) == 'odd' ? (
            <CardMedia
              controls={true}
              src={post.Url}
              component='video'
              style={{
                width: '800px',
                height: '600px',
              }}
            />
          ) : (
            <CardMedia
              src={post.Url}
              component='img'
              style={{
                width: '800px',
                height: '600px',
              }}
            />
          )}
          {post.Caption != '' && (
            <CardContent>
              <Typography variant='body2'>{post.Caption}</Typography>
            </CardContent>
          )}
          {randomMode && (
            <CardActions>
              <Button
                variant='contained'
                style={{ backgroundColor: '#EC407A', color: 'white' }}
                fullWidth={true}
                onClick={() => {
                  if (setRefresh) setRefresh(true)
                }}
              >
                Pass
              </Button>
              <Button
                variant='contained'
                style={{ backgroundColor: '#66BB6A', color: 'white' }}
                fullWidth={true}
                onClick={(e) => {
                  handleSaveBtnClick(e)
                }}
              >
                {saving ? <CircularProgress /> : 'Save'}
              </Button>
              <Menu
                open={openSaveMenu}
                anchorEl={saveMenuAnchor}
                onClose={() => {
                  setSaveMenuAnchor(null)
                }}
              >
                <Box
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    alignContent: 'center',
                    gap: '10px',
                    padding: '10px',
                  }}
                >
                  <Typography variant='subtitle1'>Caption</Typography>
                  <TextField
                    variant='outlined'
                    autoFocus={true}
                    value={savingCaption}
                    size='small'
                    onChange={(e) => {
                      setSavingCaption(e.target.value)
                    }}
                    onKeyUp={(e) => {
                      if (e.key == 'Enter') {
                        handleSavePost()
                      }
                    }}
                  />
                  <Button
                    variant='contained'
                    onClick={() => {
                      handleSavePost()
                    }}
                  >
                    <Check />
                  </Button>
                </Box>
              </Menu>
            </CardActions>
          )}
          {!randomMode && (
            <>
              <Box
                sx={{
                  margin: '5px 10px 10px',
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                }}
              >
                <div>
                  <IconButton
                    onClick={() =>
                      handleUpvotePost(
                        post,
                        currentUserState,
                        !isVoted(post.Id, currentUserState)
                      )
                    }
                  >
                    <ArrowUpward
                      color={
                        isVoted(post.Id, currentUserState)
                          ? 'primary'
                          : 'success'
                      }
                    />
                  </IconButton>
                  <IconButton onClick={() => setShowComments(!showComments)}>
                    <Chat />
                  </IconButton>
                </div>
                <Typography variant='body2'>{post.Date}</Typography>
              </Box>
              <Comments
                currentUserState={currentUserState}
                post={post}
                showComments={showComments}
              />
            </>
          )}
        </Card>
      </Grid>
    </Grid>
  )
}

export function filetype(src: string) {
  if (src) return src.substring(src.lastIndexOf('.') + 1)
}

export default Post
