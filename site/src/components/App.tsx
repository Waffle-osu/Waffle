import { useState, useEffect }  from 'react';

import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Content from './Content';
import MainLayout from './Main';

import { AppState, LoginDetails } from "./../AppState";

function App() {
	let [getLoginState, setLoginState] = useState<LoginDetails>()

	let appState: AppState = new AppState(getLoginState, setLoginState);

	useEffect(() => {
		let token = window.sessionStorage.getItem("waffle-token")

		console.log("token:" + token)

		if(token !== null) {
			let currentLoginDetails: LoginDetails = {
				username: token,
				token: token,
				userId: 0,
				loggedIn: true
			};

			appState.setLoginState(currentLoginDetails);
		}
	}, [])

	return (
		<>	
			<Router>
				<Routes>
					<Route path='/' element={<MainLayout appState={appState}></MainLayout>}>
						<Route path="beatmaps" element={<Content appState={appState}></Content>}/>
					</Route>
				</Routes>
			</Router>	
		</>
	);
}

export default App;
