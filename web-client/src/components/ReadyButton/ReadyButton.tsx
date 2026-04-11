import type { ReadyButtonProps } from "./types"
import styles from "./ReadyButton.module.css"

function ReadyButton({ isReady = false, onToggle }: ReadyButtonProps) {
    return (
        <button
            onClick={onToggle}
            className={`${styles.button} ${isReady ? styles.cancel : styles.ready}`}
        >
            {isReady ? "Cancel" : "Ready"}
        </button>    )
}

export default ReadyButton
