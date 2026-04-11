import banner from "../assets/banner.jpg"

import ReadyButton from "../components/ReadyButton/ReadyButton"
import LobbyHeader from "../components/LobbyHeader/LobbyHeader"
import { useState } from "react";

function LobbyPage() {
  const [ready, setReady] = useState(false)

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
          <LobbyHeader />

          <div
            style={{
              display: "flex",
              justifyContent: "flex-end",
              marginTop: "10px",
            }}
          >
            <ReadyButton
              isReady={ready}
              onToggle={() => setReady(!ready)}
            />
          </div>
        </div>
      </section>
    </div>
  )

}

export default LobbyPage
