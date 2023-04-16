import { Button, Container, TextField, Typography } from "@mui/material"
import "../styles/Login.css"

function Login() {
    return <Container className="main-container">
        <Typography variant="h2"> Login </Typography>
        <div className="input-container">
            <TextField label="email" variant="standard" />
            <TextField label="password" variant="standard" />
        </div>
        <Button className="login-btn" variant="contained">Login</Button>
        <Button className="sign-up-btn" style={{ textTransform: 'none' }}>Create a new account</Button>
    </Container>
}

export default Login