import React, { useEffect, useState } from 'react'
import '../styles/Root.css'
import { AppBar, Box, Tab, Tabs } from '@mui/material'
import Post from "../components/Post";
import Profile from "../components/Profile";
import { useNavigate } from 'react-router-dom';
import TopPosts from '../components/TopPosts';
import { getRandomPost } from '../utils/serverFunctions';
import Search from "../components/Search";

const tabValue = {
  top_posts: 0,
  random_posts: 1,
  search: 2,
  profile: 3
}

function Main({ tab }: { tab: "top_posts" | "random_posts" | "search" | "profile" }) {
  let navigate = useNavigate()

  let [loading, setLoading] = useState(true)
  let [response, setResponse] = useState<any>({})

  let [tabIndex, setTabIndex] = React.useState(tabValue[tab]);
  let [refreshPost, setRefreshPost] = useState(true);

  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    if (newValue == 0 || newValue == 1) {
      setRefreshPost(true)
    }
    let linksToNavigate = ["/top", "/", "/search", "/profile"]
    navigate(linksToNavigate[newValue])

    setTabIndex(newValue);
  };

  useEffect(() => {
    if (refreshPost) {
      if (tabIndex == 0) {         // top
        setLoading(true)
        getRandomPost().then(res => {
          setResponse(res)
          setLoading(false)
          setRefreshPost(false)
        })
      } else if (tabIndex == 1) {  // random
        setLoading(true)
        getRandomPost().then(res => {
          setResponse(res)
          setLoading(false)
          setRefreshPost(false)
        })
      }
    }
  }, [refreshPost])

  return (
    <>
      <AppBar position="sticky">
        <Tabs
          value={tabIndex}
          onChange={handleTabChange}
          indicatorColor="secondary"
          textColor="inherit"
          variant="fullWidth"
        >
          <Tab label="Top Posts" />
          <Tab label="Random Posts" />
          <Tab label="Search" />
          <Tab label="Profile" />
        </Tabs>
      </AppBar>
      <TabPanel value={tabIndex} index={0}>
        <TopPosts />
      </TabPanel>
      <TabPanel value={tabIndex} index={1}>
        <Post loading={loading} caption="" src={response.Result} setRefresh={setRefreshPost} random={true}
              comments={false} />
      </TabPanel>
      <TabPanel value={tabIndex} index={2}>
        <Search />
      </TabPanel>
      <TabPanel value={tabIndex} index={3}>
        <Profile />
      </TabPanel>
    </>
  )
}

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          {children}
        </Box>
      )}
    </div>
  );
}

export default Main
