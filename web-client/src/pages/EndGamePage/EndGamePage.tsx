import Button from "../../components/Button/Button";
import MainCard from "../../components/MainCard/MainCard";
import type { EndGamePageProps } from "./types"

function EndGamePage({ winner, player }: EndGamePageProps) {
  const isWinner = winner === player;

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


        <Button
          value="Retry"
          type="ready"
          handleOnClick={() => console.log("Clicked retry")}
        />
      </div>
    </MainCard>
  )
}

export default EndGamePage
