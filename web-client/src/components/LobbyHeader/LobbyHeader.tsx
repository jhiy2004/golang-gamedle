import { useState } from "react"
import RoomSpan from "../RoomSpan/RoomSpan"
import styles from "./LobbyHeader.module.css"

function LobbyHeader() {
    const [readyPlayers, setReadyPlayers] = useState(0)
    const [minPlayers, setMinPlayers] = useState(0)
    const [maxPlayers, setMaxPlayers] = useState(0)
    const [currPlayers, setCurrPlayers] = useState(0)

    const readyPlayersRatio = `Ready: ${readyPlayers} / ${currPlayers}`;

    return (
        <div className={styles.lobbyHeader}>
            <RoomSpan text={readyPlayersRatio} variant="ready" />
            <RoomSpan text={`Min: ${minPlayers}`} variant="limits" />
            <RoomSpan text={`Max: ${maxPlayers}`} variant="limits" />
        </div>
    )
}

export default LobbyHeader
