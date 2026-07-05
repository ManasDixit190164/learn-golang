import React from 'react'
import { Routes, Route, Link } from 'react-router-dom'
import Signup from './components/Signup'
import Login from './components/Login'
import Dashboard from './components/Dashboard'
import URLDetails from './components/URLDetails'
import PrivateRoute from './components/PrivateRoute'

export default function App(){
  return (
    <div className="app">
      <nav>
        <Link to="/">Home</Link>
        <Link to="/signup">Signup</Link>
        <Link to="/login">Login</Link>
      </nav>

      <main>
        <Routes>
          <Route path="/" element={<h2>Welcome to URL Shortener</h2>} />
          <Route path="/signup" element={<Signup />} />
          <Route path="/login" element={<Login />} />
          <Route path="/dashboard" element={<PrivateRoute><Dashboard/></PrivateRoute>} />
          <Route path="/urls/:id" element={<PrivateRoute><URLDetails/></PrivateRoute>} />
        </Routes>
      </main>
    </div>
  )
}
