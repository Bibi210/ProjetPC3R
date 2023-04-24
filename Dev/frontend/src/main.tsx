import React from 'react'
import ReactDOM from 'react-dom/client'
import {
    createBrowserRouter,
    RouterProvider
} from 'react-router-dom'
import Base from './views/Base'
import Login from './components/Login'
import Profile from "./components/Profile";
import Logout from "./components/Logout";

const router = createBrowserRouter([
    {path: "/", element: <Base/>},
    {path: "/login", element: <Login/>},
    {path: "/logout", element: <Logout/>},
    {path: "/profile", element: <Profile/>},
])

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
    <RouterProvider router={router}/>
)
