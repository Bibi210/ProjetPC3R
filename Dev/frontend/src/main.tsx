import ReactDOM from 'react-dom/client'
import {
    createBrowserRouter,
    RouterProvider
} from 'react-router-dom'
import Base from './views/Base'
import Login from './views/Login'
import Profile from "./components/Profile";
import Logout from "./views/Logout";

const router = createBrowserRouter([
    {path: "/", element: <Base tab="random_posts"/>},
    {path: "/top", element: <Base tab="top_posts"/>},
    {path: "/random", element: <Base tab="random_posts"/>},
    {path: "/search", element: <Base tab="search"/>},
    {path: "/profile", element: <Base tab="profile"/>},
    {path: "/login", element: <Login/>},
    {path: "/logout", element: <Logout/>},
])

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
    <RouterProvider router={router}/>
)
