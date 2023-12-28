import React, { useEffect } from 'react';
import { render, screen } from '@testing-library/react';
import TodoWrapper from '../TodoWrapper';
import { BrowserRouter } from 'react-router-dom';
import AuthProvider from '../../context/AuthProvider';
import useAuth from '../../hooks/useAuth';

const MockedTodoWrapper = () => {
    const { auth, setAuth } = useAuth();

    useEffect(() => {
        setAuth({
            id: 1,
        });
    }, []);

    return (
        auth && <TodoWrapper />
    );
}

const MockedAuthProvider = () => {
    return (
        <BrowserRouter>
            <AuthProvider>
                <MockedTodoWrapper />
            </AuthProvider>
        </BrowserRouter>
    );
}

describe('TodoWrapper', () => {
    it('should render a todo list', async () => {
        render(<MockedAuthProvider />);

        const todoEl1 = await screen.findByText("Test task 1");
        expect(todoEl1).toBeInTheDocument();

        const todoEl2 = await screen.findByText("Test task 1");
        expect(todoEl2).toBeInTheDocument();
    });
});
