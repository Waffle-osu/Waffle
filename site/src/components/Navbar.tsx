import "./Navbar.css"

import {
	BrowserRouter as Router,
	Link
} from "react-router-dom"

function Navbar() {
    return (
		<Router>		
        	<div className="navbar">
				<a href="/">Home</a>
				<a href="/download">Download</a>
				<a href="/beatmaps">Beatmaps</a>
				<a href="/leaderboards">Leaderboards</a>
				<a href="/discord">Discord</a>

				<a href="/discord" className="right-align">Register</a>
        	</div>
		</Router>
    );
}

export default Navbar;