import React from 'react'
import { BrowserRouter, Route, Routes, Navigate } from "react-router";
import Navbar from '../component/Navbar';
import Login from '../pages/auth/login/Login';
import Register from '../pages/auth/register/Register';
import Error from '../pages/error/Error';
import AuthCheck from '../guard/AuthCheck';
import Home from '../pages/home/Home';
import Auth from '../pages/auth/Auth';

const logout = () => {
  localStorage.clear();
  window.location.href = "/auth/login";
};

function MainLayout() {
  const token = localStorage.getItem("token");

  return (
    <BrowserRouter>
      {/* navbar */}
      <Navbar logout={logout} />
      <Routes>
        <Route path="/" element={
          <AuthCheck>
            <Home />
          </AuthCheck>
        } />
        <Route path="auth" element={token ? <Navigate to="/" /> : <Auth />}>
          <Route path="login" element={token ? <Navigate to="/" /> : <Login />} />
          <Route path="register" element={token ? <Navigate to="/" /> : <Register />} />
        </Route>

        {/* not found routes */}
        <Route path="*" element={<Error />} />
      </Routes>

      {/* footer component */}
    </BrowserRouter>
  );
}

export default MainLayout;
