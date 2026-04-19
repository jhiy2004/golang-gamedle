import { useEffect, useRef, useState } from "react";
import EndGamePage from "./pages/EndGamePage/EndGamePage"
import GamePage from "./pages/GamePage/GamePage"
import LobbyPage from "./pages/LobbyPage/LobbyPage"
import { createCancelMsg, createCancelRetryMsg, createGuessMsg, createReadyMsg, createRetryMsg, type Message } from "./game/message";
import { useParams, useSearchParams } from "react-router-dom";

function App() {
  const { id } = useParams()
  const [searchParams] = useSearchParams();

  const wsRef = useRef<WebSocket | null>(null);

  const [ready, setReady] = useState(false)
  const [retry, setRetry] = useState(false)

  const [readyPlayers, setReadyPlayers] = useState<number>(0)
  const [retryPlayers, setRetryPlayers] = useState<number>(0)

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


  function handleRetryClick() {
    if (!wsRef.current) return;

    const websocket = wsRef.current

    let msg
    if (retry) {
      msg = createCancelRetryMsg()
    } else {
      msg = createRetryMsg()
    }
    websocket.send(JSON.stringify(msg))
    setRetry(!retry)
  }

  function handleReadyClick() {
    if (!wsRef.current) return;
    const websocket = wsRef.current

    let msg
    if (ready) {
      msg = createCancelMsg()
    } else {
      msg = createReadyMsg()
    }
    websocket.send(JSON.stringify(msg))

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
        console.log(message.payload.playerId);
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
      case 'postGameLobby':
        console.log(message)
        setCurrPlayers(message.payload.currPlayers);
        setRetryPlayers(message.payload.retryPlayers);
        break;
      case 'restart':
        console.log(message)
        setReadyPlayers(0)
        setRetryPlayers(0)
        setState('waiting')
        setWinner('')
        setReady(false)
        setRetry(false)
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

    const playerId = searchParams.get("playerId")
    console.log("Received query param: " + playerId)

    const wsUri = `ws://${import.meta.env.VITE_APP_URL}/ws?roomId=${id}` + (playerId === null ? "" : `&playerId=${playerId}`);
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
  }, [searchParams, id]);

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
      retry={retry}
      retryPlayers={retryPlayers}
      currPlayers={currPlayers}
      handleRetryClick={handleRetryClick}
    />
  }
}

export default App
