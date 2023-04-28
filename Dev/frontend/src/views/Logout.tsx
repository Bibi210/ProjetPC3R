import { Button, Container, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import { useNavigate } from 'react-router-dom'

async function logoutServer() {
    let req = await fetch(window.location.origin + `/api/logout`, {
        method: "GET",
    })
    return await req.json()
}

function Logout() {
    const [logout, setLogout] = useState(false)
    const [callingServer, setCallingServer] = useState(true)
    const [error, setError] = useState("")
    const navigate = useNavigate()

    useEffect(() => {
        logoutServer().then((res) => {
            setCallingServer(false)
            if (res.Success == true) {
                setLogout(true)
                setTimeout(() => navigate("/"), 2000)
            } else {
                setError(res.Message)
            }
        })
    }, [])
    return <Container className="main-container">
        {callingServer && !logout &&
            <Typography variant="h2">Logging out</Typography>
        }
        {!callingServer && logout ?
            <Typography variant="h2">Logout successful</Typography>
            :
            <>
                <Typography variant="h2">Error while logging out</Typography>
                <Button
                    fullWidth
                    variant="contained"
                    style={{ backgroundColor: "#EF5350" }}
                >{error}</Button>
            </>
        }
    </Container>
}

export default Logout
