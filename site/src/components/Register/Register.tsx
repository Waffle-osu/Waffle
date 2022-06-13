import axios from "axios";
import { useState } from "react";
import { useQuery } from "react-query";
import { AppProps } from "../../AppState";

import "../Common/Content.css"
import "./Register.css"

interface EmbeddedRegisterProps {
    setVisibleState: React.Dispatch<React.SetStateAction<boolean>> | null
}

interface RegisterResponse {
    WaffleUsername: string;
    WaffleToken: string;
    WaffleUserId: number;
    WaffleRegisterStatus: string;
}

function Register(props: AppProps | EmbeddedRegisterProps) {
    const [username, setUsername] = useState<string>();
    const [password, setPassword] = useState<string>();
    const [statusText, setStatusText] = useState<string>("");

    //login query
    let {refetch: postRegisterRequest} = useQuery<RegisterResponse, Error>("login-response", async () => {
        let loginFormData = new FormData();

        loginFormData.append("username", username!)
        loginFormData.append("password", password!)

        const response = await axios.postForm<RegisterResponse>("http://127.0.0.1:80/api/waffle-site-register", loginFormData)

        return response.data
    }, {
        refetchOnWindowFocus: false,
        enabled: false
    })

    let onUsernameChange = function(change: React.FormEvent<HTMLInputElement>) {
        setUsername(change.currentTarget.value)
    }

    let onPasswordChange = function(change: React.FormEvent<HTMLInputElement>) {
        setPassword(change.currentTarget.value)
    }

    let submitRegister = async function(event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault()

        let appProps = props as AppProps

        postRegisterRequest().then((response) => {
            let responseData = response.data

            if(responseData?.WaffleToken !== "" && responseData?.WaffleUserId! > 0 && responseData?.WaffleUsername !== "") {
                window.sessionStorage.setItem("waffle-username", responseData?.WaffleUsername!)
                window.sessionStorage.setItem("waffle-token", responseData?.WaffleToken!)
                window.sessionStorage.setItem("waffle-userId", responseData?.WaffleUserId.toString()!)

                setStatusText(responseData?.WaffleRegisterStatus!)
    
                appProps.appState.setLoginState({
                    loggedIn: true,
                    userId: responseData?.WaffleUserId!,
                    username: responseData?.WaffleUsername!,
                    token: responseData?.WaffleToken!
                })

                let registerProps = (props as EmbeddedRegisterProps)

                if (registerProps.setVisibleState !== undefined) {
                    setTimeout(() => {
                        window.scrollTo({behavior:"smooth", left: 0, top: 0})

                        setTimeout(() => {
                            if(registerProps.setVisibleState !== null) {
                                registerProps.setVisibleState(false)
                            }    
                        }, 500)
                    }, 1000)         
                } 
            }
        })         
    }

    return (
        <>
            <div className="downward-content-box">
                <div className="content-item">
                    <p className="register-text">Register to Waffle!</p>

                    <div className="register-items-box">  
                        <div className="rules-box">
                            <p>Let's establish some basic rules:</p>

                            <ul>
                                <li>Don't spread this to <span className="red-text">anyone.</span> When I say anyone, I mean it. <br/>
                                    Not the client, not the link to the website, nothing.
                                </li>
                                <li>Instead of abusing glitches/bugs, report them to <span className="blue-text">Furball</span>.</li>
                                <li>No multi-accounting, <br/> 
                                    <span className="blue-text">Furball</span> keeps track of registered users heavily. <br/> 
                                    If there are more users than expected, <span className="blue-text">Furball</span> will be <span className="red-text">angry.</span>
                                </li>
                                <li>Don't cheat, there's actually no point with this small group of people.</li>
                                <li>Make sure <span className="blue-text">Furball</span> knows who you are, by this I mostly mean usernames, <br/>
                                    If you want to deviate from your normal username, let <span className="blue-text">Furball</span> know.
                                </li>
                            </ul>

                            It might be a bit extreme but I wan't to have absolute control over who gets to play on this server.
                        </div>

                        <br/>

                        <form onSubmit={submitRegister}>
                            Username:<br/> <input type="text" value={username} onChange={onUsernameChange} className="input-box"></input>

                            <br/>

                            Password:<br/> <input type="password" value={password} onChange={onPasswordChange} className="input-box"></input>
                            {statusText === "" ? (<><input type="submit" value="Register"/>  </>) : (<><p>{statusText}</p></>)}          
                        </form>                        
                    </div>
                </div>
            </div>
        </>
    )
}

export default Register;