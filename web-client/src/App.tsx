import { useEffect, useRef, useState } from "react";
import EndGamePage from "./pages/EndGamePage/EndGamePage"
import GamePage from "./pages/GamePage/GamePage"
import LobbyPage from "./pages/LobbyPage/LobbyPage"
import { createCancelMsg, createGuessMsg, createReadyMsg, type Message } from "./game/message";
import { useParams } from "react-router-dom";

function App() {
  const { id } = useParams()

  const wsRef = useRef<WebSocket | null>(null);
  const [ready, setReady] = useState(false)

  const [readyPlayers, setReadyPlayers] = useState<number>(0)
  const [minPlayers, setMinPlayers] = useState<number>(0)
  const [maxPlayers, setMaxPlayers] = useState<number>(0)
  const [currPlayers, setCurrPlayers] = useState<number>(0)
  const [players, setPlayers] = useState<string[]>([])

  // Current client player name
  const [player, setPlayer] = useState<string>('')

  const [question, setQuestion] = useState<string>('')
  const [winner, setWinner] = useState<string>('')
  const [state, setState] = useState<string>('')
  const [answer, setAnswer] = useState<string>('')
  const [playersStatus, setPlayersStatus] = useState<Record<string, number>>({})

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

    const msg = createGuessMsg(answer)
    websocket.send(JSON.stringify(msg))

    setAnswer('')
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
        setPlayer(message.payload.playerName);
        break;
      case 'notify':
        console.log(message.payload.text);
        break;
      case 'state':
        setQuestion(message.payload.question);
        setWinner(message.payload.winner);
        setState(message.payload.state);
        setPlayers(message.payload.players);
        break;
      case 'playerStatus':
        setPlayersStatus(prev => ({
          ...prev,
          [message.payload.player]: message.payload.progress
        }));
        break;
      case 'guessResponse':
        alert(message.payload.text)
        break;
      default:
        console.log("Unhandled message");
    }
  }

  useEffect(() => {
    console.log(playersStatus)
  }, [playersStatus])

  useEffect(() => {
    if (wsRef.current) return;

    const wsUri = `ws://${import.meta.env.VITE_APP_URL}?id=${id}`;
    const websocket = new WebSocket(wsUri);

    wsRef.current = websocket;

    websocket.addEventListener("error", (event) => {
      console.log("WebSocket error: ", event);
    });

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
  }, [id]);

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
      answer={answer}
      setAnswer={setAnswer}
      handleAnswerSend={handleAnswerSend}
      players={players}
      playersStatus={playersStatus}
      player={player}
    />
  } else if (state === 'end') {
    return <EndGamePage
      winner={winner}
      player={player}
    />
  }
}

export default App
