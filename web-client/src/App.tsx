import { useEffect, useRef, useState } from "react";
import EndGamePage from "./pages/EndGamePage/EndGamePage"
import GamePage from "./pages/GamePage/GamePage"
import LobbyPage from "./pages/LobbyPage/LobbyPage"
import { createCancelMsg, createGuessMsg, createReadyMsg, type Message } from "./game/message";

function App() {
  const wsRef = useRef<WebSocket | null>(null);
  const [ready, setReady] = useState(false)

  const [readyPlayers, setReadyPlayers] = useState(0)
  const [minPlayers, setMinPlayers] = useState(0)
  const [maxPlayers, setMaxPlayers] = useState(0)
  const [currPlayers, setCurrPlayers] = useState(0)
  const [question, setQuestion] = useState('')
  const [winner, setWinner] = useState('')
  const [state, setState] = useState('')
  const [answer, setAnswer] = useState('')

  function handleReadyClick() {
    if (!wsRef.current) return;
    const websocket = wsRef.current

    if (ready) {
      const msg = createCancelMsg()
      websocket.send(JSON.stringify(msg))
    } else {
      const msg = createReadyMsg()
      websocket.send(JSON.stringify(msg))
    }

    setReady(!ready)
  }

  function handleAnswerSend() {
    if (!wsRef.current) {
      return
    }

    const websocket = wsRef.current

    if (ready) {
      const msg = createGuessMsg(answer)
      websocket.send(JSON.stringify(msg))
    } else {
      const msg = createReadyMsg()
      websocket.send(JSON.stringify(msg))
    }

    setReady(!ready)
  }

  function handleServerResponse(message: Message) {
    switch (message.cmd) {
      case 'lobby':
        setCurrPlayers(message.payload.currPlayers);
        setReadyPlayers(message.payload.readyPlayers);
        break;
      case 'start':
        setMinPlayers(message.payload.minPlayers);
        setMaxPlayers(message.payload.maxPlayers);
        break;
      case 'notify':
        console.log(message.payload.text);
        break;
      case 'state':
        setQuestion(message.payload.question);
        setWinner(message.payload.winner);
        setState(message.payload.state);
        break;

      default:
        console.log("Unhandled message");
    }
  }

  useEffect(() => {
    if (wsRef.current) return;

    const wsUri = `ws://${import.meta.env.VITE_APP_URL}`;
    const websocket = new WebSocket(wsUri);

    wsRef.current = websocket;

    websocket.onopen = () => {
      console.log("Connected");
    };

    websocket.onmessage = (e) => {
      console.log("Received:", e.data);
      const message: Message = JSON.parse(e.data)
      handleServerResponse(message)
    };

    return () => {
      websocket.close();
      wsRef.current = null;
    };
  }, []);

  if (state === 'waiting' || state === '') {
    return <LobbyPage
        ready={ready}
        handleReadyClick={handleReadyClick}
        readyPlayers={readyPlayers}
        minPlayers={minPlayers}
        maxPlayers={maxPlayers}
        currPlayers={currPlayers}
      />
  } else if(state === 'playing') {
    return <GamePage
      question={question}
      setAnswer={setAnswer}
      handleAnswerSend={handleAnswerSend}
    />
  } else if (state === 'end') {
    return <EndGamePage/>
  }
}

export default App
