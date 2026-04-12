import type { GamePageProps } from "./types"
import styles from "../../components/ReadyButton/ReadyButton.module.css"

function GamePage({ question, setAnswer, handleAnswerSend }: GamePageProps) {
  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        minHeight: "100vh",
        backgroundColor: "#f3f4f6",
      }}
    >
      <section
        style={{
          width: 320,
          borderRadius: "16px",
          overflow: "hidden",

          backgroundColor: "white",
          boxShadow: "0 10px 25px rgba(0, 0, 0, 0.15)",

          paddingBottom: "12px",
        }}
      >
        <div style={{ padding: "10px" }}>
          <div>
            <h1>Question</h1>
            <p>{question}</p>
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
                onChange={(e) => setAnswer(e.target.value)}
              />
            </div>
            <div
              style={{
                display: "flex",
                justifyContent: "flex-end",
                marginTop: "10px",
              }}
            >
              <button className={`${styles.button} ${styles.ready}`} onClick={handleAnswerSend}>Send Answer</button>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}

export default GamePage
