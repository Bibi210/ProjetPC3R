import ReactDOM from 'react-dom/client'
import { createTheme, ThemeOptions } from '@mui/material/styles'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import Main from './views/Main'
import Login from './views/Login'
import Logout from './views/Logout'
import { ThemeProvider } from '@emotion/react'

const router = createBrowserRouter([
  { path: '/', element: <Main tab='random_posts' /> },
  { path: '/top', element: <Main tab='top_posts' /> },
  { path: '/random', element: <Main tab='random_posts' /> },
  { path: '/search', element: <Main tab='search' /> },
  { path: '/profile', element: <Main tab='profile' /> },
  { path: '/login', element: <Login /> },
  { path: '/logout', element: <Logout /> },
])

export const theme: ThemeOptions = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#673ab7',
    },
    secondary: {
      main: '#e91e63',
    },
    error: {
      main: '#c62828',
    },
    warning: {
      main: '#ffcc80',
    },
    info: {
      main: '#616161',
    },
    success: {
      main: '#81c784',
    },
  },
  shape: {
    borderRadius: 10,
  },
  spacing: 8,
})

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <ThemeProvider theme={theme}>
    <RouterProvider router={router} />
  </ThemeProvider>
)
