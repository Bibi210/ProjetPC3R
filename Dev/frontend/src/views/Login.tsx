import {
  Button,
  Card,
  CardContent,
  CircularProgress,
  Collapse,
  Container,
  TextField,
  Typography,
} from '@mui/material'
import '../styles/Login.css'
import { useRef, useState } from 'react'
import { Navigate } from 'react-router-dom'
import { Notification, NotificationType } from '../utils/types'
import { createAccount, login as loginAccount } from '../utils/serverFunctions'

function Login() {
  const [createAccountMode, setCreateAccountMode] = useState(false)

  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')
  const [password2, setPassword2] = useState('')

  const [notifications, setNotifications] = useState<Notification[]>([])
  const [loggedIn, setLoggedIn] = useState(false)
  const [sendingRequest, setSendingRequest] = useState(false)

  let notificationsRef = useRef<Notification[]>(notifications)
  notificationsRef.current = notifications

  function addNotif(msg: string, type: NotificationType) {
    function unique_id(id: number): number {
      for (const a of notifications) {
        if (a.id == id) {
          return unique_id(id + 1)
        }
      }
      return id
    }

    let id = unique_id(notifications.length - 1)
    let notification: Notification = { id, msg, type, show: true }
    let newNotifsState = [...notifications, notification]
    setNotifications(newNotifsState)
    setTimeout(() => {
      setNotifications(
        notificationsRef.current.map((n) => {
          if (n.id == id) {
            n.show = false
          }
          return n
        })
      )
    }, 4500)
    setTimeout(() => {
      setNotifications(notificationsRef.current.filter((n) => n.id != id))
    }, 5000)
  }

  function validateBeforeRequest(login: string, pass: string) {
    if (login == '') {
      addNotif('please add a username', NotificationType.ERROR)
      return
    }
    if (pass == '') {
      addNotif('please add a password', NotificationType.ERROR)
      return
    }
    if (!createAccountMode) {
      setSendingRequest(true)
      loginAccount(login, pass).then((res) => {
        setSendingRequest(false)
        if (res.Success) {
          setLoggedIn(true)
        } else {
          addNotif(res.Message, NotificationType.ERROR)
        }
      })
      return loginAccount(login, pass)
    }
    if (password2 != password) {
      addNotif("passwords don't match", NotificationType.ERROR)
      return
    }
    setSendingRequest(true)
    createAccount(login, pass).then((res) => {
      setSendingRequest(false)
      if (res.Success) {
        addNotif('Successfully created account', NotificationType.INFO)
        setCreateAccountMode(false)
      } else {
        addNotif(res.Message, NotificationType.ERROR)
      }
    })
  }

  return (
    <Container className='main-container'>
      {loggedIn && <Navigate to='/' />}
      <Card>
        <CardContent style={{ padding: '50px' }}>
          <Typography variant='h2' color='text.primary' marginBottom='30px'>
            {' '}
            {createAccountMode ? 'Create an account' : 'Login'}{' '}
          </Typography>
          <div className='errors'>
            {notifications.map((n) => (
              <Collapse appear={true} key={n.id} in={n.show}>
                <Button
                  fullWidth
                  variant='contained'
                  onClick={() =>
                    setNotifications(
                      notifications.filter((notif) => notif.id != n.id)
                    )
                  }
                  style={{
                    backgroundColor:
                      n.type == NotificationType.ERROR ? '#EF5350' : '#3F51B5',
                    marginBottom: '10px',
                  }}
                >
                  {n.msg}
                </Button>
              </Collapse>
            ))}
          </div>
          <div className='input-container' style={{marginBottom: "50px"}}>
            <TextField
              label='username'
              variant='filled'
              error={login == ''}
              helperText={login == '' ? 'username cannot be empty' : ''}
              onChange={(e) => setLogin(e.currentTarget.value)}
            />
            <TextField
              label='password'
              type='password'
              variant='filled'
              error={password == ''}
              helperText={password == '' ? 'password cannot be empty' : ''}
              onChange={(e) => setPassword(e.currentTarget.value)}
            />
            {createAccountMode && (
              <TextField
                label='retype password'
                type='password'
                variant='filled'
                error={password2 == ''}
                helperText={
                  password2 == '' ? 'please re enter your password' : ''
                }
                onChange={(e) => setPassword2(e.currentTarget.value)}
              />
            )}
          </div>
          <Button
            className='login-btn'
            variant='contained'
            onClick={() => validateBeforeRequest(login, password)}
          >
            {sendingRequest ? (
              <CircularProgress />
            ) : createAccountMode ? (
              'Create account'
            ) : (
              'Login'
            )}
          </Button>
          <Button
            className='sign-up-btn'
            style={{ textTransform: 'none' }}
            onClick={() => setCreateAccountMode(!createAccountMode)}
          >
            {createAccountMode
              ? 'Already have an account? Login'
              : 'Create a new account'}
          </Button>
        </CardContent>
      </Card>
    </Container>
  )
}

export default Login
