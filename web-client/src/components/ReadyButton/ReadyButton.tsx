import type { ReadyButtonProps } from "./types"
import Button from "../Button/Button"

function ReadyButton({ isReady = false, onToggle }: ReadyButtonProps) {
    return (
        <Button
            value={ isReady ? "Cancel" : "Ready" }
            type={ isReady ? "cancel" : "ready" }
            handleOnClick={onToggle}
        />
    )
}

export default ReadyButton
