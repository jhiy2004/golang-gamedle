import type { MainCardProps } from "./types"
import styles from "./MainCard.module.css"

function MainCard({ children }: MainCardProps) {
    return (
        <div className={styles.mainCardBackground}>
            <section className={styles.mainCard}>
                { children }
            </section>
        </div>
    )
}

export default MainCard
