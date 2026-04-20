import Button from "../../components/Button/Button"
import MainCard from "../../components/MainCard/MainCard"
import type { GamePageProps } from "./types"

function GamePage({ question, answer, setAnswer, handleAnswerSend, players, playersStatus, player }: GamePageProps) {
  const sortedPlayers = [
    player,
    ...players.filter(p => p !== player)
  ]

  return (
    <MainCard>
      <div>
        <div style={{ padding: "10px" }}>
          <div>
            <h1>Question</h1>
            <p style={{maxWidth: "320px"}}>{question}</p>
          </div>
          <div>
            <div>
              <input
                style={{
                  padding: "12px 14px",
                  borderRadius: "10px",
                  border: "1px solid #ccc",
                  fontSize: "1rem",
                  outline: "none",
                  transition: "all 0.2s ease",
                  width: "100%",
                }}
                type="text"
                value={answer}
                onChange={(e) => setAnswer(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    handleAnswerSend()
                  }
                }}
                autoFocus
              />
            </div>
            <div
              style={{
                display: "flex",
                justifyContent: "flex-end",
                marginTop: "10px",
              }}
            >
              <Button
                value="Send Answer"
                type="ready"
                handleOnClick={handleAnswerSend}
              />
            </div>
          </div>

          <div>
            {sortedPlayers.map((p) => {
              return <p key={p}>{(player !== p) ? p : `${p} (You)`}  {'*'.repeat(playersStatus[p] ?? 0)}</p>
            })}
          </div>
        </div>
      </div>
    </MainCard>
  )
}

export default GamePage
