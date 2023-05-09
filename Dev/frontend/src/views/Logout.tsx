import { Button, Container, Typography } from '@mui/material'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { logout } from '../utils/serverFunctions'

function Logout() {
  const [loggedOut, setLoggedOut] = useState(false)
  const [callingServer, setCallingServer] = useState(true)
  const [error, setError] = useState('')
  const navigate = useNavigate()

  useEffect(() => {
    logout().then((res) => {
      setCallingServer(false)
      if (res.Success) {
        setLoggedOut(true)
        setTimeout(() => navigate('/'), 2000)
      } else {
        setError(res.Message)
      }
    })
  }, [])
  return (
    <Container className='main-container'>
      {callingServer && !loggedOut && (
        <Typography variant='h2' color='text.primary'>
          Logging out
        </Typography>
      )}
      {!callingServer && loggedOut ? (
        <Typography variant='h2' color='text.primary'>
          Logout successful
        </Typography>
      ) : (
        <>
          <Typography variant='h2' color='text.primary'>
            Error while logging out
          </Typography>
          <Button
            fullWidth
            variant='contained'
            style={{
              marginTop: "50px"
            }}
          >
            {error}
          </Button>
        </>
      )}
    </Container>
  )
}

export default Logout
