import "./Navbar.css"

import { Link } from "react-router-dom"

import { AppProps } from "./../AppState";

function Navbar(props: AppProps) {
	let userPanel: JSX.Element;

	if(props.appState.loginState?.loggedIn) {
		userPanel = (
			<>
				<Link to={"/users/" + props.appState.loginState.userId}>
					<div className="login-area">
						<p>{props.appState.loginState.username}</p>
						<img src={"http://127.0.0.1:80/a/" + props.appState.loginState.userId } ></img>
						
					</div>
				</Link>
			</>
		)
	} else {
		userPanel = (
			<>
				<Link to="/login">Log In</Link>
				<Link to="/register">Register</Link>
			</>
		);
	}

    return (	
        <div className="navbar">
			<p className="waffle-text">Waffle</p>

			<Link to="/">Home</Link>
			<Link to="/download">Download</Link>
			<Link to="/beatmaps">Beatmaps</Link>
			<Link to="/leaderboards">Leaderboards</Link>
			<Link to="/discord">Discord</Link>

			<div className="right-align">
				{userPanel}
			</div>
        </div>
    );
}

export default Navbar;