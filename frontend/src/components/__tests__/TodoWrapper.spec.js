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
    it('should have a valid form', async () => {
        render(<MockedAuthProvider />);
        const formEl = await screen.findByTestId(/todo-item-1/i);
        expect(formEl).toBeInTheDocument();
    });
});
