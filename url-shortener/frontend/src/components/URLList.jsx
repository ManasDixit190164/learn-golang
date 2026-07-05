import React, { useEffect, useState } from 'react'
import api from '../api'
import { Link } from 'react-router-dom'

export default function URLList(){
  const [urls, setUrls] = useState([])
  const [err, setErr] = useState(null)

  const load = async () => {
    try{
      const res = await api.get('/urls')
      setUrls(res.data.data)
    }catch(e){ setErr('Failed to load') }
  }

  useEffect(()=>{ load() },[])

  if (err) return <div>{err}</div>
  return (
    <div>
      <ul>
        {urls.map(u=> (
          <li key={u.id}>
            <a href={u.short_url} target="_blank" rel="noreferrer">{u.short_url}</a>
            {' - '}{u.original_url}
            {' - '}
            <Link to={`/urls/${u.id}`}>Analytics</Link>
          </li>
        ))}
      </ul>
    </div>
  )
}
