import React, { useState } from 'react'
import api from '../api'

export default function CreateURL(){
  const [form, setForm] = useState({original_url:'', custom_alias:'', title:''})
  const [msg, setMsg] = useState(null)

  const submit = async e => {
    e.preventDefault()
    try{
      const body = { original_url: form.original_url }
      if (form.custom_alias) body.custom_alias = form.custom_alias
      if (form.title) body.title = form.title
      const res = await api.post('/urls', body)
      setMsg('Created: ' + res.data.data.short_url)
      setForm({original_url:'', custom_alias:'', title:''})
    }catch(err){ setMsg(err.response?.data?.error || 'Failed') }
  }

  return (
    <form onSubmit={submit} className="panel">
      <h4>Create Short URL</h4>
      {msg && <div className="info">{msg}</div>}
      <input aria-label="Original URL" placeholder="Original URL" required type="url" value={form.original_url} onChange={e=>setForm({...form,original_url:e.target.value})} />
      <input aria-label="Custom alias" placeholder="Custom alias (optional)" value={form.custom_alias} onChange={e=>setForm({...form,custom_alias:e.target.value})} />
      <input aria-label="Title" placeholder="Title (optional)" value={form.title} onChange={e=>setForm({...form,title:e.target.value})} />
      <button type="submit" className="btn full">Create</button>
    </form>
  )
}
