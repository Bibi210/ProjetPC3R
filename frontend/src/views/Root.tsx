import { useEffect, useState } from 'react'
import '../styles/Root.css'
import { Container, Typography } from '@mui/material'


async function callServerRoot() {
  let res = await fetch("http://localhost:4242/")
  return (await res.text())
}

function Root() {
  let [state, setState] = useState("loading")
  useEffect(() => {
    callServerRoot().then(res => setState(res))
  }, [])
  return (
    <Container>
      {state == "loading" ?
        <Typography variant='h3'>Calling server on /</Typography> :
        (<>
          <Typography variant='h3'>Response from Server:</Typography>
          <Typography>{state}</Typography>
        </>)
      }
    </Container>
  )
}

export default Root
