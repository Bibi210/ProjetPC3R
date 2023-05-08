import React, { useEffect, useState } from 'react'
import '../styles/Root.css'
import { AppBar, Box, Grid, Tab, Tabs } from '@mui/material'
import Post from '../components/Post'
import Profile from '../components/Profile'
import { useNavigate } from 'react-router-dom'
import TopPosts from '../components/TopPosts'
import {
  getPrivateProfile,
  getPublicProfile,
  getRandomPost,
} from '../utils/serverFunctions'
import Search from '../components/Search'
import { User, ServerResponse } from '../utils/types'

const tabValue = {
  top_posts: 0,
  random_posts: 1,
  search: 2,
  profile: 3,
}

function Main({
  tab,
}: {
  tab: 'top_posts' | 'random_posts' | 'search' | 'profile'
}) {
  const navigate = useNavigate()

  const [loading, setLoading] = useState(true)
  const [response, setResponse] = useState<ServerResponse<string> | null>(null)

  const [tabIndex, setTabIndex] = React.useState(tabValue[tab])
  const [refreshPost, setRefreshPost] = useState(true)
  const [currentUser, setCurrentUser] = useState<User | null>(null)

  function refreshCurrentUser() {
    getPrivateProfile().then((res) => {
      if (res.Success) {
        if (res.Result) {
          let privateUser: User = res.Result
          getPublicProfile(privateUser.Username).then((resPublic) => {
            if (resPublic.Success) {
              setCurrentUser(resPublic.Result)
            }
          })
        } else {
          console.error(res.Message)
          alert(res.Message)
        }
      } else {
        console.error(res.Message)
        alert(res.Message)
      }
    })
  }

  useEffect(() => {
    refreshCurrentUser()
  }, [])

  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    if (newValue == 0 || newValue == 1) {
      setRefreshPost(true)
    }
    let linksToNavigate = ['/top', '/', '/search', '/profile']
    navigate(linksToNavigate[newValue])

    setTabIndex(newValue)
  }

  useEffect(() => {
    if (refreshPost) {
      if (tabIndex == 0) {
        // top
        setLoading(true)
        getRandomPost().then((res) => {
          setResponse(res)
          setLoading(false)
          setRefreshPost(false)
        })
      } else if (tabIndex == 1) {
        // random
        setLoading(true)
        getRandomPost().then((res) => {
          setResponse(res)
          setLoading(false)
          setRefreshPost(false)
        })
      }
    }
  }, [refreshPost])

  return (
    <div style={{backgroundColor: "black", minHeight: "100vh"}}>
      <AppBar position='sticky'>
        <Tabs
          value={tabIndex}
          onChange={handleTabChange}
          indicatorColor='secondary'
          textColor='inherit'
          variant='fullWidth'
        >
          <Tab label='Top Posts' />
          <Tab label='Random Posts' />
          <Tab label='Search' />
          <Tab label='Profile' />
        </Tabs>
      </AppBar>
      <TabPanel value={tabIndex} index={0}>
        <TopPosts
          currentUserState={{
            get: currentUser,
            set: setCurrentUser,
            refresh: refreshCurrentUser,
          }}
        />
      </TabPanel>
      <TabPanel value={tabIndex} index={1}>
        <Grid container justifyContent='center' marginBottom={4}>
          <Grid item>
            <Post
              currentUserState={{
                get: currentUser,
                set: setCurrentUser,
                refresh: refreshCurrentUser,
              }}
              loading={loading}
              setRefresh={setRefreshPost}
              randomMode={true}
              post={{
                Id: -1,
                Caption: '',
                Creator: '',
                Date: '',
                Upvotes: 0,
                Url: response ? response.Result : '',
                CommentIds: [],
              }}
            />
          </Grid>
        </Grid>
      </TabPanel>
      <TabPanel value={tabIndex} index={2}>
        <Search
          currentUserState={{
            get: currentUser,
            set: setCurrentUser,
            refresh: refreshCurrentUser,
          }}
        />
      </TabPanel>
      <TabPanel value={tabIndex} index={3}>
        <Profile
          currentUserState={{
            get: currentUser,
            set: setCurrentUser,
            refresh: refreshCurrentUser,
          }}
        />
      </TabPanel>
    </div>
  )
}

type TabPanelProps = {
  children?: React.ReactNode
  index: number
  value: number
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props

  return (
    <div
      role='tabpanel'
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  )
}

export default Main
