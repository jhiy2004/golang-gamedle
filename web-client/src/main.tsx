import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import './index.css'

import {
  createBrowserRouter,
  RouterProvider
} from 'react-router-dom'
import EnterRoomPage from './pages/EnterRoomPage/EnterRoomPage'

const router = createBrowserRouter([
  {
    path: "/enter",
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
