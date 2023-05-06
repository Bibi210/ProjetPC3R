import ReactDOM from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import Main from './views/Main'
import Login from './views/Login'
import Logout from "./views/Logout";

const router = createBrowserRouter([
  { path: "/", element: <Main tab="random_posts" /> },
  { path: "/top", element: <Main tab="top_posts" /> },
  { path: "/random", element: <Main tab="random_posts" /> },
  { path: "/search", element: <Main tab="search" /> },
  { path: "/profile", element: <Main tab="profile" /> },
  { path: "/login", element: <Login /> },
  { path: "/logout", element: <Logout /> },
])

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <RouterProvider router={router} />
)
