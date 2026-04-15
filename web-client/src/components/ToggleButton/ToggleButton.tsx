import type { ToggleButtonProps } from "./types"
import Button from "../Button/Button"

function ToggleButton({ isToggled = false, onToggle, toggleText }: ToggleButtonProps) {
    return (
        <Button
            value={ isToggled ? "Cancel" : toggleText }
            type={ isToggled ? "cancel" : "ready" }
            handleOnClick={onToggle}
        />
    )
}

export default ToggleButton
