import React, {useEffect, useState} from 'react'
import '../styles/Root.css'
import {AppBar, Box, Tab, Tabs, Typography} from '@mui/material'
import Post from "../components/Post";
import Profile from "../components/Profile";


async function callServerRandomShitPost() {
    let res = await fetch("http://localhost:25565/api/random_shitpost")
    return await res.json()
}

function Base() {
    let [loading, setLoading] = useState(true)
    let [response, setResponse] = useState<any>({})

    let [tab, setTab] = React.useState(0);
    let [refreshPost, setRefreshPost] = useState(true);

    const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
        setRefreshPost(true)
        setTab(newValue);
    };

    useEffect(() => {
        if (refreshPost) {
            if (tab == 0) {         // top
                setLoading(true)
                callServerRandomShitPost().then(res => {
                    setResponse(res)
                    setLoading(false)
                    setRefreshPost(false)
                })
            } else if (tab == 1) {  // random
                setLoading(true)
                callServerRandomShitPost().then(res => {
                    setResponse(res)
                    setLoading(false)
                    setRefreshPost(false)
                })
            }
        }
    }, [refreshPost])

    return (
        <>
            <AppBar position="static">
                <Tabs
                    value={tab}
                    onChange={handleTabChange}
                    indicatorColor="secondary"
                    textColor="inherit"
                    variant="fullWidth"
                >
                    <Tab label="Top Posts"/>
                    <Tab label="Random Posts"/>
                    <Tab label="Search"/>
                    <Tab label="Conversations"/>
                    <Tab label="Profile"/>
                </Tabs>
            </AppBar>
            <TabPanel value={tab} index={0}>
                <Post loading={loading} src={response.Result} setRefresh={setRefreshPost}/>
            </TabPanel>
            <TabPanel value={tab} index={1}>
                <Post loading={loading} src={response.Result} setRefresh={setRefreshPost}/>
            </TabPanel>
            <TabPanel value={tab} index={2}>
                Search
            </TabPanel>
            <TabPanel value={tab} index={3}>
                Conversations
            </TabPanel>
            <TabPanel value={tab} index={4}>
                <Profile/>
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
    const {children, value, index, ...other} = props;

    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`simple-tabpanel-${index}`}
            aria-labelledby={`simple-tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box sx={{p: 3}}>
                    <Typography>{children}</Typography>
                </Box>
            )}
        </div>
    );
}

export default Base
