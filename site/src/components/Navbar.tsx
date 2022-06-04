import "./Navbar.css"

import { Link } from "react-router-dom"

import { AppProps } from "./../AppState";

function Navbar(props: AppProps) {
    return (	
        <div className="navbar">
			<p className="waffle-text">Waffle</p>

			<Link to="/">Home</Link>
			<Link to="/download">Download</Link>
			<Link to="/beatmaps">Beatmaps</Link>
			<Link to="/leaderboards">Leaderboards</Link>
			<Link to="/discord">Discord</Link>

			<div className="right-align">
				<p>{props.appState.loginState?.loggedIn ? props.appState.loginState?.username : "Log In"}</p>
			</div>
        </div>
    );
}

export default Navbar;