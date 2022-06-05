import { useEffect } from "react";

import { AppProps } from "../../AppState";

import "./../Common/Content.css"

function Beatmaps(props: AppProps) {
    return (
        <>
            <div className="downward-content-box">
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