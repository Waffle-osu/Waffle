import { useEffect } from "react";

import { AppProps } from "../../AppState";

import "./../Common/Content.css"
import BeatmapElement from "./BeatmapElement/BeatmapElement";
import BeatmapListBox from "./BeatmapListBox/BeatmapListBox";
import BeatmapSearchBox from "./BeatmapSearchBox/BeatmapSearchBox";

function Beatmaps(props: AppProps) {
    let newQuerySubmitHandler = (query: string) => {

    }

    return (
        <>
            <div className="downward-content-box">
                <BeatmapSearchBox appState={props} onNewQuerySubmit={newQuerySubmitHandler}></BeatmapSearchBox>

                <BeatmapListBox></BeatmapListBox>
            </div>
        </>
    );
}

export default Beatmaps;