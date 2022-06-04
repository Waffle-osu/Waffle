import { useState } from "react";
import { AppProps } from "../AppState";

import "./Content.css"
import "./Login.css"

function Login(props: AppProps) {
    const [username, setUsername] = useState<string>();
    const [password, setpassword] = useState<string>();

    let onUsernameChange = function(change: React.FormEvent<HTMLInputElement>) {
        setUsername(change.currentTarget.value)
    }

    let onPasswordChange = function(change: React.FormEvent<HTMLInputElement>) {
        setpassword(change.currentTarget.value)
    }

    let submitLogin = function(event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault()

        console.log("Logging in " + username + " with Password: " + password)
    }

    return (
        <>
            <div className="downward-content-box">
                <div className="content-item ">
                    <p className="login-text">Log in</p>

                    <div className="login-items-box">   
                        
                        <form onSubmit={submitLogin}>
                            Username:<br/> <input type="text" value={username} onChange={onUsernameChange} className="input-box"></input>
                        
                            <br/>
                        
                            Password:<br/> <input type="password" value={password} onChange={onPasswordChange} className="input-box"></input>

                            <input type="submit" value="Log In"/>
                        </form>
                    </div>
                </div>
            </div>
        </>
    )
}

export default Login;