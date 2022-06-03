import React from 'react';

import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Content from './Content';
import MainLayout from './Main';

function App() {
	return (
		<>	
			<Router>
				<Routes>
					<Route path='/' element={<MainLayout></MainLayout>}>
						<Route path="beatmaps" element={<Content></Content>}/>
					</Route>
				</Routes>
			</Router>
			
		</>
	);
}

export default App;
