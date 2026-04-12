import styles from "./Button.module.css"
import type { ButtonProps, ButtonType } from "./types"

function Button({ value, type, handleOnClick }: ButtonProps) {
    function typeToStyle(type: ButtonType) {
        switch (type) {
            case 'ready':
                return styles.ready
            case 'cancel':
                return styles.cancel
            default:
                return styles.ready
        }
    }

    return (
        <button
            onClick={handleOnClick}
            className={`${styles.button} ${typeToStyle(type)}`}
        >
            { value }
        </button>
    )
}

export default Button
