export type EndGamePageProps = {
    winner: string;
    player: string;
    retry: boolean;
    retryPlayers: number;
    currPlayers: number;
    handleRetryClick: () => void
}
