export type ButtonProps = {
    value: string;
    type: ButtonType;
    handleOnClick: () => void;
}

export type ButtonType = "ready" | "cancel";
