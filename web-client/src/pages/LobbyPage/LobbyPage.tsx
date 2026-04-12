import banner from "../../assets/banner.jpg"

import ReadyButton from "../../components/ReadyButton/ReadyButton"
import LobbyHeader from "../../components/LobbyHeader/LobbyHeader"
import type { LobbyPageProps } from "./types"
import MainCard from "../../components/MainCard/MainCard"

function LobbyPage({
  ready,
  handleReadyClick,
  readyPlayers,
  minPlayers,
  maxPlayers,
  currPlayers,
}: LobbyPageProps) {
  return (
    <MainCard>
      <div>
        <img
          src={banner}
          alt="Banner"
          style={{
            width: "100%",
            height: 140,
            objectFit: "cover",
          }}
        />

        <div style={{ padding: "10px" }}>
          <LobbyHeader 
            readyPlayers={readyPlayers}
            minPlayers={minPlayers}
            maxPlayers={maxPlayers}
            currPlayers={currPlayers}
          />

          <div
            style={{
              display: "flex",
              justifyContent: "flex-end",
              marginTop: "10px",
            }}
          >
            <ReadyButton
              isReady={ready}
              onToggle={handleReadyClick}
            />
          </div>
        </div>
      </div>
    </MainCard>
  )

}

export default LobbyPage
