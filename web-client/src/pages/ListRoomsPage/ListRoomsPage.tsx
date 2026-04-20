import { useEffect, useState } from "react"
import MainCard from "../../components/MainCard/MainCard"
import type { ListRoomsResponseDTO, RoomResponseDTO } from "./types"
import { useNavigate } from "react-router-dom"
import Button from "../../components/Button/Button"

function ListRoomsPage() {
    const [rooms, setRooms] = useState<RoomResponseDTO[]>([])
    const navigate = useNavigate()

    async function listRooms() : Promise<ListRoomsResponseDTO> {
        const baseUrl = import.meta.env.VITE_APP_URL
        const response = await fetch(`http://${baseUrl}/list`, {
            method: 'GET',
        });

        if (!response.ok) {
            alert('Failed to list rooms')
        }

        const data: ListRoomsResponseDTO = await response.json()
        return data
    }

    useEffect(() => {
        async function fetchRooms() {
            const data = await listRooms()
            setRooms(data.rooms)
        }

        fetchRooms()
    }, [])

    return (
        <MainCard>
            <section
                style={{
                    padding: "1rem",
                    display: "flex",
                    flexDirection: "column",
                    gap: "0.5rem"
                }}
            >
                <h1 style={{ marginBottom: "1rem" }}>Rooms</h1>

                {rooms.length === 0 ? <p>No Rooms</p>:rooms.map((r) => (
                    <Button
                        key={r.id}
                        value={`Room\n#${r.id}`}
                        type='ready' handleOnClick={() => navigate(`/room/${r.id}`)}
                    />
                ))}
            </section>
        </MainCard>
    )
}

export default ListRoomsPage
