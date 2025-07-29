import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import './index.css'
import App from './App.tsx'
import Register from './components/Register'
import Login from './components/Login'
import { AuthProvider } from './context/authContext.tsx'
import Dashboard from './components/Dashboard.tsx'
import { NewPrint } from './components/NewPrint.tsx'
import Administrator from './components/Administrator.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<App />} />
          <Route path="/register" element={<Register />} />
          <Route path="/login" element={<Login />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/new" element={<NewPrint />} />
          <Route path="/admin" element={<Administrator/>} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  </StrictMode>,
)