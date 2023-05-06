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
        <Typography variant='h2'>Logging out</Typography>
      )}
      {!callingServer && loggedOut ? (
        <Typography variant='h2'>Logout successful</Typography>
      ) : (
        <>
          <Typography variant='h2'>Error while logging out</Typography>
          <Button
            fullWidth
            variant='contained'
            style={{ backgroundColor: '#EF5350' }}
          >
            {error}
          </Button>
        </>
      )}
    </Container>
  )
}

export default Logout
