import {Container, Typography} from "@mui/material";
import {useEffect, useState} from "react";

async function logoutServer() {
    let req = await fetch(window.location.origin + `/api/logout`, {
        method: "PUT",
    })
    return await req.json()
}

function Logout() {
    let [logout, setLogout] = useState(false)
    let [callingServer, setCallingServer] = useState(true)
    let [error, setError] = useState("")

    useEffect(() => {
        logoutServer().then((res) => {
            setCallingServer(false)
            if (res.Success == true) {
                setLogout(true)
            } else {
                setError(res.Message)
            }
        })
    }, [])
    return <Container className="main-container">
        {callingServer && !logout &&
            <Typography variant="h2">Logging out</Typography>
        }
        {!callingServer && logout &&
            <Typography variant="h2">Succesfully Logged Out</Typography>
        }
        {!callingServer && !logout &&
            <>
                <Typography variant="h2">Error while logging out</Typography>
                {error}
            </>
        }
    </Container>
}

export default Logout
