import "./Navbar.css"

import {
	BrowserRouter as Router,
	Route,
	Link
} from "react-router-dom"

function Navbar() {
    return (	
        <div className="navbar">
			<Link to="/">Home</Link>
			<Link to="/download">Download</Link>
			<Link to="/beatmaps">Beatmaps</Link>
			<Link to="/leaderboards">Leaderboards</Link>
			<Link to="/discord">Discord</Link>

			<a href="/discord" className="right-align">Register</a>
        </div>
    );
}

export default Navbar;