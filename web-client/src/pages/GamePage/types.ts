export type GamePageProps = {
    question: string;
    answer: string;
    setAnswer: React.Dispatch<React.SetStateAction<string>>
    handleAnswerSend: () => void;
    players: string[];
    playersStatus: Record<string, number>;
}
