export type GamePageProps = {
    question: string;
    setAnswer: React.Dispatch<React.SetStateAction<string>>
    handleAnswerSend: () => void;
}
