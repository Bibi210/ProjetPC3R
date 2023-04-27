import { Button, CircularProgress, Container, TextField, Typography } from "@mui/material"
import "../styles/Login.css"
import { useState } from "react"
import { Navigate } from "react-router-dom";

async function request(email: string, password: string, action: string) {
    console.log(email, password)
    let req = await fetch(window.location.origin + `/api/${action}`, {
        method: action == "login" ? "PUT" : "POST",
        headers: {
            "Content-Type": "text/plain"
        },
        mode: "cors",
        body: `{"Login":"${email}", "Mdp":"${password}"}`,
    })
    let json = await req.json()
    return json
}

enum NotificationType { ERROR, NOTIF }
type Notification = {
    msg: string,
    type: NotificationType
}

function Login() {
    const [createAccountMode, setCreateAccountMode] = useState(false)

    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [password2, setPassword2] = useState("")

    const [notifications, setNotifications] = useState<Notification[]>([])
    const [loggedIn, setLoggedIn] = useState(false)
    const [sendingRequest, setSendingRequest] = useState(false)

    function addNotif(notification: Notification) {
        setNotifications([notification, ...notifications])
        setTimeout(() => {
            setNotifications(notifications.filter((err) => err != notification))
        }, 5000)
    }

    function validateBeforeRequest(email: string, pass: string, action: string) {
        if (email == "") {
            addNotif({ msg: "please add an email", type: NotificationType.ERROR })
            return
        }
        if (pass == "") {
            addNotif({ msg: "please add a password", type: NotificationType.ERROR })
            return
        }
        if (!createAccountMode) {
            setSendingRequest(true)
            request(email, pass, action).then(res => {
                setSendingRequest(false)
                if (res.Success) {
                    setLoggedIn(true)
                } else {
                    addNotif({ msg: res.Message, type: NotificationType.ERROR })
                }
            })
            return request(email, pass, action)
        }
        if (password2 != password) {
            addNotif({ msg: "passwords don't match", type: NotificationType.ERROR })
            return
        }
        setSendingRequest(true)
        request(email, pass, action).then(res => {
            setSendingRequest(false)
            if (res.Success) {
                addNotif({ msg: "account successfully created", type: NotificationType.NOTIF })
                setCreateAccountMode(false)
            } else {
                addNotif({ msg: res.Message, type: NotificationType.ERROR })
            }
        })
    }

    return <Container className="main-container">
        {loggedIn && <Navigate to="/" />}
        <Typography variant="h2"> {createAccountMode ? "Create an account" : "Login"} </Typography>
        <div className="errors">
            {notifications.map((error) =>
                <Button
                    fullWidth
                    variant="contained"
                    style={{ backgroundColor: error.type == NotificationType.ERROR ? "#EF5350" : "#3F51B5" }}
                >{error.msg}</Button>)
            }
        </div>
        <div className="input-container">
            <TextField
                label="email"
                variant="filled"
                error={email == ""}
                helperText={email == "" ? "email cannot be empty" : ""}
                onChange={(e) => setEmail(e.currentTarget.value)} />
            <TextField
                label="password"
                type="password"
                variant="filled"
                error={password == ""}
                helperText={password == "" ? "password cannot be empty" : ""}
                onChange={(e) => setPassword(e.currentTarget.value)}
            />
            {createAccountMode &&
                <TextField
                    label="retype password"
                    type="password"
                    variant="filled"
                    error={password2 == ""}
                    helperText={password2 == "" ? "please re enter your password" : ""}
                    onChange={(e) => setPassword2(e.currentTarget.value)}
                />}
        </div>
        <Button
            className="login-btn"
            variant="contained"
            onClick={() => validateBeforeRequest(email, password, createAccountMode ? "create_account" : "login")}
        >{sendingRequest ? <CircularProgress /> : createAccountMode ? "Create account" : "Login"}</Button>
        <Button
            className="sign-up-btn"
            style={{ textTransform: 'none' }}
            onClick={() => setCreateAccountMode(!createAccountMode)}
        >
            {createAccountMode ? "Already have an account? Login" : "Create a new account"}
        </Button>
    </Container>
}

export default Login
