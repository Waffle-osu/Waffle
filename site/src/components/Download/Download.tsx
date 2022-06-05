import { createRef, MouseEventHandler, useRef, useState, LegacyRef } from "react";
import { AppProps } from "../../AppState";
import Register from "../Register/Register";

import "./../Common/Content.css"
import "./Download.css"

function DownloadPanel(props: AppProps) {
    let [showRegisterPanel, setShowRegisterPanel] = useState(false);

    let registerStepClick = function(event: React.MouseEvent<HTMLElement>) {
        event.preventDefault()

        setShowRegisterPanel(!showRegisterPanel)
       
        setTimeout(() => {
             registerMoveDiv.current?.scrollIntoView({behavior: "smooth"})
        }, 50) 
    }

    let registerMoveDiv = useRef<HTMLElement>();

    return (
        <>
            <div className="downward-content-box">
                <div className="content-item">
                

                    <p className="started-text">Let's get you started!</p>

                    <div className="step-panels">  
                        <a className="right-flowing-container step-one" href="" onClick={registerStepClick}>
                            Step 1. <br/> <br/>

                            <p className="centered-text">Create Waffle Account</p>
         
                            <p className="centered-text middle-aligned-text">
                                You're required to have a Waffle account to play with others on Waffle!
                            </p>

                            <p className="down-aligned-text small-text">Click on this if you need to create one!</p>
                        </a>

                        <a className="right-flowing-container step-two" href="">
                            Step 2. <br/> <br/>

                            <p className="centered-text">Download Waffle Updater</p> 

                            <p className="centered-text middle-aligned-text">
                                This will download the Waffle osu! client and which are required to connect to Waffle
                            </p>
                        </a>

                        <a className="right-flowing-container step-three" href="">
                            Step 3. <br/> <br/>

                            <p className="centered-text">Log in and enjoy!</p>

                            <p className="centered-text middle-aligned-text">
                                If you got to this part, then you can just log into the client and enjoy! Have fun playing Waffle!
                            </p>                            
                        </a>
                    </div>

                    <br/>

                    <div className="info-text">
                        <p>
                            Once you're logged into <span className="blue-text">Waffle</span>,
                            you can use the built-in <span className="pink-text">osu!direct</span> for free
                            To more conviniently download maps, however it is also possible to download
                            Beatmaps from the <a className="dark-blue-text" href="/beatmaps">Beatmaps page</a> on the Website

                            <br/> <br/>

                            <span className="blue-text">Waffle</span> aims to recreate the osu! 2011 experience,
                            that means that there is <span className="red-text">no Performance Ranking</span>,
                            Ranked Score is also awarded on a <span className="green-text">per-mapset basis</span>,
                            and obviously a 2011 <span className="pink-text">osu! client</span> is used for gameplay.

                            All maps from 2007 to 2011 are ranked (or in some cases <span className="orange-text">Approved</span>)
                            no further maps will be added to the ranked list for period accuracy.
                        </p>
                    </div>
                </div>

                {showRegisterPanel ? (
                    <>        
                        <div ref={registerMoveDiv as LegacyRef<HTMLDivElement> | undefined}></div>
                        <Register appState={props.appState} setVisibleState={setShowRegisterPanel}></Register>
                    </>
                ) : (<></>)}
            </div>
        </>
    )
}

export default DownloadPanel