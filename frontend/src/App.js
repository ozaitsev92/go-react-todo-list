import React from 'react';
import { Routes, Route } from 'react-router-dom';
import './App.css';
import SignIn from './pages/SignIn';
import SignUp from './pages/SignUp';
import TodoList from './pages/TodoList';
import NotFoundPage from './pages/NotFoundPage';
import Layout from './components/Layout';
import RequireAuth from './components/RequireAuth';

function App() {
    return (
        <Routes>
            <Route path="/" element={<Layout />}>
                <Route path="/signin" element={<SignIn />} />
                <Route path="/signup" element={<SignUp />} />
                <Route element={<RequireAuth />}>
                    <Route path="/" element={<TodoList />} />
                </Route>
                <Route path="*" element={<NotFoundPage />} />
            </Route>
        </Routes>
    );
}

export default App;
