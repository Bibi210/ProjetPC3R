import React from 'react'
import ReactDOM from 'react-dom/client'
import {
  createBrowserRouter,
  RouterProvider
} from 'react-router-dom'
import Root from './views/Root'
import Login from './views/Login'

const router = createBrowserRouter([
  { path: "/", element: <Root /> },
  { path: "/hello", /* element: */ },
  { path: "/login", element: <Login /> },
  { path: "/get_profile" },
  { path: "/create_account" },
  { path: "/logout" },
  { path: "/post/:postId"}
])

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
)
