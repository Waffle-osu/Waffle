import React from 'react';
import './App.css';
import Navbar from './Navbar';

import { BrowserRouter as Router, Routes, Route } from "react-router-dom";

function App() {
	return (
		<div className="main">
			<Navbar></Navbar>


			<Router>
				<Routes>
					<Route path='beatmaps' element={
						<p>Beatmaps!</p>
					}/>
				</Routes>
			</Router>
		</div>
	);
}

export default App;
