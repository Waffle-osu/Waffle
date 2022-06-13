import { KeyboardEvent, LegacyRef, useRef } from "react";
import { AppProps } from "../../../AppState"

import "./BeatmapSearchBox.css"

interface BeatmapSearchBoxProps {
    appState: AppProps;
    onNewQuerySubmit (query: string): void;
}

function BeatmapSearchBox(props: BeatmapSearchBoxProps) {
    let queryRef = useRef<HTMLInputElement>()

    let mostPlayedOnClick = (event: React.MouseEvent<HTMLElement>) => {
        event.preventDefault()

        queryRef.current!.value! = "Most Played"

        querySubmitHandler(queryRef.current!.value!)
    }

    let newestOnClick = (event: React.MouseEvent<HTMLElement>) => {
        event.preventDefault()

        queryRef.current!.value! = "Newest"

        querySubmitHandler(queryRef.current!.value!)
    }

    let topRatedOnClick = (event: React.MouseEvent<HTMLElement>) => {
        event.preventDefault()

        queryRef.current!.value! = "Top Rated"

        querySubmitHandler(queryRef.current!.value!)
    }

    let querySubmitHandler = (query: string) => {
        props.onNewQuerySubmit(query)
    }

    let queryOnKeyPressedHandler = (event: KeyboardEvent) => {
        if (event.key === 'Enter') {
            querySubmitHandler(queryRef.current!.value!)
        }        
    }    

    return (
        <>
            <div className="content-item">
                <p className="search-header-text">Beatmap Search</p>

                <div className="quick-search-button-container">
                    <a className="quick-search-button search-reset-button" href="" onClick={mostPlayedOnClick}>
                        <p>Most Played</p>
                    </a>

                    <a className="quick-search-button search-newest-button" href="" onClick={newestOnClick}>
                        <p>Newest Maps</p>
                    </a>

                    <a className="quick-search-button search-top-rated-button" href="" onClick={topRatedOnClick}>
                        <p>Top Rated</p>
                    </a>
                </div>

                <p className="search-header-text">Search Query</p>

                <input type="text" className="search-query-input-box" ref={queryRef as LegacyRef<HTMLInputElement> | undefined} onKeyDown={queryOnKeyPressedHandler}></input>
            </div>
        </>
    )
}

export default BeatmapSearchBox