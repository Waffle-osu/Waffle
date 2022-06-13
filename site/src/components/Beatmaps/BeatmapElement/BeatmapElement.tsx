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
            rankedImage = (<><img src="/selection-ranked.png" height="24" className="ranked-status"></img></>)
            break
        case 2:
            rankedImage = (<><img src="/selection-approved.png" height="24" className="ranked-status"></img></>)
            break 
    }

    let thumbnailUrl = "http://127.0.0.1:80/mt/" + props.BeatmapSetId + "l"

    let elementOnClick = (event: React.MouseEvent<HTMLElement>) => {
        event.preventDefault()
    }

    return (
        <>
            <div className="beatmap-element-container">
                <a className="beatmap-element-anchor" href="" onClick={elementOnClick}>
                    <div className="beatmap-element">    
                        <div className="left-side-container">
                            <div className="thumbnail">
                                <img src={thumbnailUrl} height="48" className="beatmap-thumbnail"></img>
                            </div>
                            <div className="metadata">
                                <p className="beatmap-metadata">
                                    {props.TitleString}
                                </p>
                            </div>
                            
                            <div className="extra-metadata">
                                {props.LastUpdatedString} 
                            </div>
                            
                            <div className="beatmap-rating-container">
                                <progress className="beatmap-rating" max="100" value={props.MapRating * 10.0}></progress> 
                            </div>
                        </div>

                        <div className="right-side-container">
                            <h4>{props.Creator}</h4>

                            {rankedImage}
                        </div>
                    </div>
                </a>
            </div>
        </>
    )
}

export default BeatmapElement;