import RoomSpan from "../RoomSpan/RoomSpan"
import styles from "./LobbyHeader.module.css"
import type { LobbyHeaderProps } from "./types";

function LobbyHeader({
  readyPlayers,
  minPlayers,
  maxPlayers,
  currPlayers,
}: LobbyHeaderProps) {

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
