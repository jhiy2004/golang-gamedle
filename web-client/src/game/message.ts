export type PostGameLobbyMsg = {
    cmd: 'postGameLobby';
    payload: {
        retryPlayers: number;
        currPlayers: number;
    }
}

export type GuessResponseMsg = {
    cmd: 'guessResponse';
    payload: {
        correct: boolean;
        text: string;
    }
}

export type PlayerStatusMsg = {
    cmd: 'playerStatus';
    payload: {
        player: string;
        progress: number;
    }
}

export type NotifyMsg = {
    cmd: 'notify';
    payload: {
        text: string;
    }
}

export type StateMsg = {
    cmd: 'state';
    payload: {
        question: string;
        winner: string;
        state: string;
        players: string[];
    }
}

export type GuessMsg = {
    cmd: 'guess';
    payload: {
        answer: string
    }
}

export type LobbyMsg = {
    cmd: 'lobby';
    payload: {
        currPlayers: number;
        readyPlayers: number;
    }
}

export type StartMsg = {
    cmd: 'start';
    payload: {
        minPlayers: number;
        maxPlayers: number;
        playerName: string;
        playerId: string;
    }
}

export type ReadyMsg = {
    cmd: 'ready';
    // eslint-disable-next-line @typescript-eslint/no-empty-object-type
    payload: {}
}

export type CancelMsg = {
    cmd: 'cancel';
    // eslint-disable-next-line @typescript-eslint/no-empty-object-type
    payload: {}
}

export type RetryMsg = {
    cmd: 'retry';
    // eslint-disable-next-line @typescript-eslint/no-empty-object-type
    payload: {}
}

export type CancelRetryMsg = {
    cmd: 'cancelRetry';
    // eslint-disable-next-line @typescript-eslint/no-empty-object-type
    payload: {}
}

export type RestartMsg = {
    cmd: 'restart';
    // eslint-disable-next-line @typescript-eslint/no-empty-object-type
    payload: {}
}

export type Message =
    | NotifyMsg
    | StateMsg
    | GuessMsg
    | LobbyMsg
    | StartMsg
    | ReadyMsg
    | CancelMsg
    | PlayerStatusMsg
    | GuessResponseMsg
    | PostGameLobbyMsg
    | RetryMsg
    | CancelRetryMsg
    | RestartMsg

export function createLobbyMsg(currPlayers: number, readyPlayers: number): LobbyMsg {
    const lobbyMsg: LobbyMsg = {
        cmd: 'lobby',
        payload: {
            currPlayers,
            readyPlayers,
        }
    } 
    
    return lobbyMsg
}


export function createNotifyMsg(text: string): NotifyMsg {
    const notifyMsg: NotifyMsg = {
        cmd: 'notify',
        payload: {
            text,
        }
    } 
    
    return notifyMsg
}

export function createStateMsg(question: string, winner: string, state: string, players: string[]): StateMsg {
    const stateMsg: StateMsg = {
        cmd: 'state',
        payload: {
            question,
            winner,
            state,
            players,
        }
    } 
    
    return stateMsg
}


export function createGuessMsg(answer: string): GuessMsg {
    const guessMsg: GuessMsg = {
        cmd: 'guess',
        payload: {
            answer
        }
    } 
    
    return guessMsg
}

export function createStartMsg(minPlayers: number, maxPlayers: number, playerName: string, playerId: string): StartMsg {
    const startMsg: StartMsg = {
        cmd: 'start',
        payload: {
            minPlayers,
            maxPlayers,
            playerName,
            playerId,
        }
    } 
    
    return startMsg
}

export function createReadyMsg(): ReadyMsg {
    const readyMsg: ReadyMsg = {
        cmd: 'ready',
        payload: { }
    }
    
    return readyMsg
}

export function createCancelMsg(): CancelMsg {
    const cancelMsg: CancelMsg = {
        cmd: 'cancel',
        payload: {}
    }
    
    return cancelMsg
}


export function createRetryMsg(): RetryMsg {
    const retryMsg: RetryMsg = {
        cmd: 'retry',
        payload: { }
    }
    
    return retryMsg
}

export function createCancelRetryMsg(): CancelRetryMsg {
    const cancelRetryMsg: CancelRetryMsg = {
        cmd: 'cancelRetry',
        payload: {}
    }
    
    return cancelRetryMsg
}
