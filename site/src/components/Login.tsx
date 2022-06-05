import axios from "axios";
import { useState } from "react";
import { useQuery } from "react-query";
import { useNavigate } from "react-router-dom";
import { AppProps } from "../AppState";

import "./Content.css"
import "./Login.css"

interface LoginResponse {
    WaffleUsername: string;
    WaffleToken: string;
    WaffleUserId: number;
}

function Login(props: AppProps) {
    const [username, setUsername] = useState<string>();
    const [password, setPassword] = useState<string>();

    //login query
    let {status: loginStatus, data: loginDataResponse, refetch: postLoginRequest} = useQuery<LoginResponse, Error>("login-response", async () => {
        let loginFormData = new FormData();

        loginFormData.append("username", username!)
        loginFormData.append("password", password!)

        const response = await axios.postForm<LoginResponse>("http://127.0.0.1:80/api/waffle-login", loginFormData)

        return response.data
    }, {
        refetchOnWindowFocus: false,
        enabled: false
    })

    let navigate = useNavigate()

    let onUsernameChange = function(change: React.FormEvent<HTMLInputElement>) {
        setUsername(change.currentTarget.value)
    }

    let onPasswordChange = function(change: React.FormEvent<HTMLInputElement>) {
        setPassword(change.currentTarget.value)
    }

    let submitLogin = async function(event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault()

        await postLoginRequest().then((response) => {
            const resposneData = response.data

            if(resposneData?.WaffleToken !== "" && resposneData?.WaffleUserId! > 0 && resposneData?.WaffleUsername !== "") {
                window.sessionStorage.setItem("waffle-username", resposneData?.WaffleUsername!)
                window.sessionStorage.setItem("waffle-token", resposneData?.WaffleToken!)
                window.sessionStorage.setItem("waffle-userId", resposneData?.WaffleUserId.toString()!)
    
                props.appState.setLoginState({
                    loggedIn: true,
                    userId: resposneData?.WaffleUserId!,
                    username: resposneData?.WaffleUsername!,
                    token: resposneData?.WaffleToken!
                })
                
                navigate("/")
            }
        })
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

                            <p>{
                                loginStatus === "error" ? (<>Login error occured!</>) : 
                                loginDataResponse?.WaffleUserId! > 0 ? (<>Wrong Username and or Password!</>) : (<></>)
                            }</p>
                        </form>
                    </div>
                </div>
            </div>
        </>
    )
}

export default Login;