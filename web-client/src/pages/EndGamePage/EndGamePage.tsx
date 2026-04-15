import MainCard from "../../components/MainCard/MainCard";
import ToggleButton from "../../components/ToggleButton/ToggleButton";
import RoomSpan from "../../components/RoomSpan/RoomSpan";
import type { EndGamePageProps } from "./types"

function EndGamePage({ winner, player, retry, retryPlayers, currPlayers, handleRetryClick}: EndGamePageProps) {
  const isWinner = winner === player;

  const retryPlayersRatio = `Ready: ${retryPlayers} / ${currPlayers}`;

  return (
    <MainCard>
      <div style={{textAlign: "center", padding: "24px"}}>
        <h1 style={{ fontSize: "1.5rem", margin: 0 }}>
          {isWinner
            ? "🎉 You won the game!"
            : `🏆 ${winner} won the game`}
        </h1>

        <p style={{ color: "#555", margin: 0 }}>
          {isWinner
            ? "Great job! You answered the questions correctly."
            : "Better luck next time — try again!"}
        </p>

        <div
          style={{
            display: "flex",
            justifyContent: "center",
            marginTop: "10px",
          }}
        >
          <RoomSpan text={retryPlayersRatio} variant="ready"/>
        </div>

        <div
          style={{
            display: "flex",
            justifyContent: "center",
            marginTop: "10px",
          }}
        >
          <ToggleButton
            isToggled={retry}
            onToggle={handleRetryClick}
            toggleText="Retry"
          />
        </div>
      </div>
    </MainCard>
  )
}

export default EndGamePage
