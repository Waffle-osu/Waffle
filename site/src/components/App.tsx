import { useState, useEffect }  from 'react';

import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Beatmaps from './Beatmaps/Beatmaps';
import MainLayout from './Main';

import { AppState, LoginDetails } from "./../AppState";
import Login from './Login/Login';
import { QueryClient, QueryClientProvider } from 'react-query';
import DownloadPanel from './Download/Download';
import Register from './Register/Register';

const queryClient = new QueryClient();

function App() {
	let [getLoginState, setLoginState] = useState<LoginDetails>()

	let appState: AppState = new AppState(getLoginState, setLoginState);

	useEffect(() => {
		let username = window.sessionStorage.getItem("waffle-username")
		let token = window.sessionStorage.getItem("waffle-token")
		let userId = window.sessionStorage.getItem("waffle-userId")

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
			<QueryClientProvider  client={queryClient}>
				<Router>
					<Routes>
						<Route path='/' element={ <MainLayout appState={appState} ></MainLayout> }>
							<Route path="beatmaps" element={ <Beatmaps appState={appState} ></Beatmaps> }/>
							<Route path="login" element={ <Login appState={appState} ></Login> }/>
							<Route path="download" element={ <DownloadPanel appState={appState} ></DownloadPanel> }/>
							<Route path="register" element={ <Register appState={appState} ></Register> }/>
						</Route>
					</Routes>
				</Router>	
			</QueryClientProvider>
		</>
	);
}

export default App;
