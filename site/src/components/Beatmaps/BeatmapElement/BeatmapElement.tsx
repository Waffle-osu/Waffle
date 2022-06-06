import "./BeatmapElement.css"

interface BeatmapElementProps {
    TitleString: string;
    LastUpdatedString: string;
    MapRating: number;
    Creator: string;
    BeatmapSetId: string
    RankedStatus: number;
}

function BeatmapElement(props: BeatmapElementProps) {
    let rankedImage: JSX.Element = (<></>);

    switch(props.RankedStatus) {
        case 0:
            rankedImage = (<></>)
            break
        case 1:
            rankedImage = (<><img src="/selection-ranked.png" height="32" className="ranked-status"></img></>)
            break
        case 2:
            rankedImage = (<><img src="/selection-approved.png" height="32" className="ranked-status"></img></>)
            break 
    }

    return (
        <>
            <div className="beatmap-element-container">
                <a className="beatmap-element-anchor" href="">
                    <div className="beatmap-element">    
                        <div className="left-side-container">
                            <div className="thumbnail">
                                <img src="http://127.0.0.1:80/mt/22472l" height="64" className="beatmap-thumbnail"></img>
                            </div>
                            <div className="metadata">
                                <p className="beatmap-metadata">
                                    {props.TitleString}
                                </p>
                            </div>
                            
                            <div className="extra-metadata">
                                {props.LastUpdatedString} <br/>
                                Plays: 123 <br/>
                            </div>
                            
                            <progress className="beatmap-rating" max="100" value={props.MapRating * 10.0}></progress>
                        </div>

                        <div className="right-side-container">
                            <h3>{props.Creator}</h3>

                            {rankedImage}
                        </div>
                    </div>
                </a>
            </div>
        </>
    )
}

export default BeatmapElement;