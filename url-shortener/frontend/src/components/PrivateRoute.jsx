import React from 'react'
import { Navigate } from 'react-router-dom'
import { getAccessToken } from '../api'

export default function PrivateRoute({ children }){
  const token = getAccessToken()
  if (!token) return <Navigate to="/login" replace />
  return children
}
