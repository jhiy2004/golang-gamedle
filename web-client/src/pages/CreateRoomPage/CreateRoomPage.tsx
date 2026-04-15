import { useNavigate } from "react-router-dom";
import Button from "../../components/Button/Button"
import MainCard from "../../components/MainCard/MainCard"

type CreateRoomResponse = {
  id: string;
}

function CreateRoomPage() {
  const navigate = useNavigate()

  async function createRoom(): Promise<CreateRoomResponse> {
    const baseUrl = import.meta.env.VITE_APP_URL

    const response = await fetch(`http://${baseUrl}/create`, {
      method: 'POST',
    });

    if (!response.ok) {
      throw new Error(`Error: ${response.status}`)
    }

    const data: CreateRoomResponse = await response.json();
    return data
  }

  async function handleCreateRoom() {
    try {
      const response = await createRoom()
      navigate(`/room/${response.id}`)
    } catch (error) {
      console.error(error)
    }
  }

  return (
    <MainCard>
      <Button
        value="Create Room"
        type="ready"
        handleOnClick={handleCreateRoom}
      />
    </MainCard>
  )
}

export default CreateRoomPage
