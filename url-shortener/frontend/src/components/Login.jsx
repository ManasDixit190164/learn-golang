import React, { useState } from 'react'
import api, { setTokens } from '../api'
import { useNavigate } from 'react-router-dom'

export default function Login(){
  const [form, setForm] = useState({email:'', password:''})
  const [error, setError] = useState(null)
  const nav = useNavigate()

  const submit = async (e) => {
    e.preventDefault()
    try{
      const res = await api.post('/auth/login', form)
      setTokens(res.data.data)
      nav('/dashboard')
    }catch(err){ setError(err.response?.data?.error || 'Failed') }
  }

  return (
    <form onSubmit={submit} className="panel">
      <h3>Login</h3>
      {error && <div className="error">{error}</div>}
      <input aria-label="Email" placeholder="Email" required type="email" value={form.email} onChange={e=>setForm({...form,email:e.target.value})} />
      <input aria-label="Password" type="password" placeholder="Password" required value={form.password} onChange={e=>setForm({...form,password:e.target.value})} />
      <button type="submit" className="btn full">Login</button>
    </form>
  )
}
