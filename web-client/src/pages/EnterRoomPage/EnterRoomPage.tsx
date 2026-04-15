import { useState } from "react"
import Button from "../../components/Button/Button"
import MainCard from "../../components/MainCard/MainCard"
import { useNavigate } from "react-router-dom"

function EnterRoomPage() {
    const [roomId, setRoomId] = useState<string>('')
    const navigate = useNavigate()

    function handleSendClick() {
        navigate(`/room/${roomId}`)
    }

    return (
        <MainCard>
            <div style={{ padding: "10px" }}>
                <input
                    style={{
                        padding: "12px 14px",
                        borderRadius: "10px",
                        border: "1px solid #ccc",
                        fontSize: "1rem",
                        outline: "none",
                        transition: "all 0.2s ease",
                        width: "100%",
                    }}
                    type="text"
                    value={roomId}
                    onChange={(e) => setRoomId(e.target.value)}
                    onKeyDown={(e) => {
                        if (e.key === 'Enter') {
                            handleSendClick()
                        }
                    }}
                    autoFocus
                />

                <div
                    style={{
                        display: "flex",
                        justifyContent: "flex-end",
                        marginTop: "10px",
                    }}
                >
                    <Button
                        value="Send Answer"
                        type="ready"
                        handleOnClick={handleSendClick}
                    />
                </div>
            </div>

        </MainCard>
    )

}

export default EnterRoomPage
