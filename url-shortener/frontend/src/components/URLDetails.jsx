import React, { useEffect, useState } from 'react'
import api from '../api'
import { useParams } from 'react-router-dom'

export default function URLDetails(){
  const { id } = useParams()
  const [data, setData] = useState(null)
  const [err, setErr] = useState(null)

  useEffect(()=>{
    const load = async ()=>{
      try{
        const res = await api.get(`/urls/${id}/analytics`)
        setData(res.data.data)
      }catch(e){ setErr('Failed') }
    }
    load()
  },[id])

  if (err) return <div>{err}</div>
  if (!data) return <div>Loading...</div>
  return (
    <div>
      <h3>Analytics</h3>
      <pre>{JSON.stringify(data,null,2)}</pre>
    </div>
  )
}
