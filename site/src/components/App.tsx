import { useState, useEffect }  from 'react';

import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Beatmaps from './Beatmaps';
import MainLayout from './Main';

import { AppState, LoginDetails } from "./../AppState";

function App() {
	let [getLoginState, setLoginState] = useState<LoginDetails>()

	let appState: AppState = new AppState(getLoginState, setLoginState);

	useEffect(() => {
		let username = window.sessionStorage.getItem("waffle-username")
		let token = window.sessionStorage.getItem("waffle-token")
		let userId = window.sessionStorage.getItem("waffle-userId")

		console.log("token:" + token)

		if(token !== null && username !== null && userId !== null) {
			let currentLoginDetails: LoginDetails = {
				username: username,
				token: token,
				userId: Number(userId),
				loggedIn: true
			};

			appState.setLoginState(currentLoginDetails);
		}
	}, [])

	return (
		<>	
			<Router>
				<Routes>
					<Route path='/' element={ <MainLayout appState={appState} ></MainLayout> }>
						<Route path="beatmaps" element={ <Beatmaps appState={appState} ></Beatmaps> }/>
					</Route>
				</Routes>
			</Router>	
		</>
	);
}

export default App;
