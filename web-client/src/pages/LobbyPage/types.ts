export type LobbyPageProps = {
    ready: boolean;
    handleReadyClick: () => void;
    readyPlayers: number;
    minPlayers: number;
    maxPlayers: number;
    currPlayers: number;
}
