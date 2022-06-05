import React from "react";

import './App.css'

import Navbar from "./Navbar/Navbar";
import { AppProps } from "./../AppState";

import { Outlet } from 'react-router-dom';

function MainLayout(props: AppProps) {
    return (
        <>
            <div className="main">
                <Navbar appState={props.appState}></Navbar>
                
                <Outlet/>
			</div>
        </>
    )
}

export default MainLayout;