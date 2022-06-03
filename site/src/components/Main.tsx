import React from "react";

import './App.css'

import Navbar from "./Navbar";

import { Outlet } from 'react-router-dom';

function MainLayout() {
    return (
        <>
            <div className="main">
				<Navbar></Navbar>
                
                <Outlet/>
			</div>
        </>
    )
}

export default MainLayout;