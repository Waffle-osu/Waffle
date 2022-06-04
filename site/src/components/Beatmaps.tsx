import { useEffect } from "react";

import { AppProps } from "../AppState";

import "./Beatmaps.css"

function Beatmaps(props: AppProps) {
    return (
        <>
            <div className="beatmaps-content">
                <div className="content-item">
                    <p>Test</p>
                    <br></br>
                    <br></br>
                    <br></br>
                    <br></br>
                </div>

                <div className="content-item">
                    <p>Test</p>
                    <br></br>
                </div>

                <div className="content-item">
                    <p>Test</p>
                    <br></br>
                    <br></br>
                    <br></br>
                    <br></br>
                </div>
            </div>
        </>
    );
}

export default Beatmaps;