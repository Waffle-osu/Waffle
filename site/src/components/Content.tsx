import { useEffect } from "react";

import { AppProps } from "./../AppState";

function Content(props: AppProps) {
    useEffect(() => {
        window.sessionStorage.setItem("waffle-token", "Furball")
		
        props.appState.setLoginState({
            username: "Furball",
            token: "dslfhsfskjdfhsdf",
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