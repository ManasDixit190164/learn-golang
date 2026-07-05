import React from 'react'
import CreateURL from './CreateURL'
import URLList from './URLList'
import { clearTokens } from '../api'
import { useNavigate } from 'react-router-dom'

export default function Dashboard(){
  const nav = useNavigate()
  const logout = () => { clearTokens(); nav('/login') }
  return (
    <div>
      <h2>Your URLs</h2>
      <button onClick={logout} className="btn secondary">Logout</button>
      <CreateURL />
      <URLList />
    </div>
  )
}
