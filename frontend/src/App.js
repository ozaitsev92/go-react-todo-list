import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import './App.css';
import SignIn from './pages/SignIn';
import SignUp from './pages/SignUp';
import TodoList from './pages/TodoList';
import NotFoundPage from './pages/NotFoundPage';

function App() {
    return (
        <main className="App">
            <BrowserRouter>
                <Routes>
                    <Route path="/" element={<TodoList />} />
                    <Route path="/signin" element={<SignIn />} />
                    <Route path="/signup" element={<SignUp />} />
                    <Route path="*" element={<NotFoundPage />} />
                </Routes>
            </BrowserRouter>
        </main>
    );
}

export default App;
