import styles from "./RoomSpan.module.css"

import type { RoomSpanProps } from "./types"

function RoomSpan({ text, variant = "limits" }: RoomSpanProps) {
    const variantClass = {
        limits: styles["badge--limits"],
        ready: styles["badge--ready"],
    }[variant]

    return (
        <span className={`${styles.badge} ${variantClass}`}>
            {text}
        </span>
    )
}

export default RoomSpan
