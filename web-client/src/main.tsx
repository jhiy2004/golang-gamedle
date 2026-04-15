import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import './index.css'

import {
  createBrowserRouter,
  RouterProvider
} from 'react-router-dom'
import EnterRoomPage from './pages/EnterRoomPage/EnterRoomPage'
import CreateRoomPage from './pages/CreateRoomPage/CreateRoomPage.tsx'

const router = createBrowserRouter([
  {
    path: "/create",
    element: <CreateRoomPage />
  },
  {
    path: "/",
    element: <EnterRoomPage />
  },
  {
    path: "/room/:id",
    element: <App />
  }
])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router={router}/>
  </StrictMode>,
)
