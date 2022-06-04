import { useEffect } from "react";

import { AppProps } from "./../AppState";

function Content(props: AppProps) {
    useEffect(() => {
        window.sessionStorage.setItem("waffle-token", "Hi!")
		
        props.appState.setLoginState({
            username: "Hi!",
            token: "Hi!",
            userId: 2,
            loggedIn: true
        })
    }, [])

    return (
        <>
            <p>Beatmaps!</p>
            <p>Beatmaps!</p>
            <p>Beatmaps!</p>
            <p>Beatmaps!</p>
        </>
    );
}

export default Content;