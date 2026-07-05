import React, { useState } from 'react'
import api, { setTokens } from '../api'
import { useNavigate } from 'react-router-dom'

export default function Signup(){
  const [form, setForm] = useState({name:'', email:'', password:''})
  const [error, setError] = useState(null)
  const nav = useNavigate()

  const submit = async (e) => {
    e.preventDefault()
    try{
      const res = await api.post('/auth/signup', form)
      setTokens(res.data.data)
      nav('/dashboard')
    }catch(err){ setError(err.response?.data?.error || 'Failed') }
  }

  return (
    <form onSubmit={submit} className="panel">
      <h3>Signup</h3>
      {error && <div className="error">{error}</div>}
      <input aria-label="Name" placeholder="Name" required type="text" value={form.name} onChange={e=>setForm({...form,name:e.target.value})} />
      <input aria-label="Email" placeholder="Email" required type="email" value={form.email} onChange={e=>setForm({...form,email:e.target.value})} />
      <input aria-label="Password" type="password" placeholder="Password" required minLength={6} value={form.password} onChange={e=>setForm({...form,password:e.target.value})} />
      <button type="submit" className="btn full">Signup</button>
    </form>
  )
}
