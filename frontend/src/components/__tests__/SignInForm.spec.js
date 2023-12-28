import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import SignInForm from '../SignInForm';
import { BrowserRouter } from 'react-router-dom';
import { act } from 'react-dom/test-utils';

const MockedSignInForm = () => {
    return (
        <BrowserRouter>
            <SignInForm />
        </BrowserRouter>
    );
}

describe('SignInForm', () => {
    beforeEach(() => {
        act(() => {
            render(<MockedSignInForm />);
        });
    });

    it('should have a valid form', () => {
        const formEl = screen.getByTestId(/signin-form/i);
        expect(formEl).toBeInTheDocument();

        const emailLabelEl = screen.getByLabelText(/Email/i);
        expect(emailLabelEl).toBeInTheDocument();

        const emailInputEl = screen.getByPlaceholderText(/Email/i);
        expect(emailInputEl).toBeInTheDocument();

        const passwordLabelEl = screen.getByLabelText(/Password/i);
        expect(passwordLabelEl).toBeInTheDocument();

        const passwordInputEl = screen.getByPlaceholderText(/Password/i);
        expect(passwordInputEl).toBeInTheDocument();

        const buttonEl = screen.getByRole("button", { name: /Sign In/i });
        expect(buttonEl).toBeInTheDocument();
    });

    it('should not sign in if the email is empty', async () => {
        const emailInputEl = screen.getByPlaceholderText(/Email/i);
        expect(emailInputEl).toBeInTheDocument();

        const passwordInputEl = screen.getByPlaceholderText(/Password/i);
        expect(passwordInputEl).toBeInTheDocument();

        fireEvent.change(emailInputEl, {target: {value: ''}})
        fireEvent.change(passwordInputEl, {target: {value: 'Zaq12wsx'}})

        await act(() => fireEvent.submit(screen.getByTestId(/signin-form/i)));

        const errorEl = screen.getByText(/Invalid email or password./i);
        expect(errorEl).toBeInTheDocument();
    });

    it('should not sign in if the password is empty', async () => {
        const emailInputEl = screen.getByPlaceholderText(/Email/i);
        expect(emailInputEl).toBeInTheDocument();

        const passwordInputEl = screen.getByPlaceholderText(/Password/i);
        expect(passwordInputEl).toBeInTheDocument();

        fireEvent.change(emailInputEl, {target: {value: 'test@example.com'}});
        fireEvent.change(passwordInputEl, {target: {value: ''}});

        await act(() => fireEvent.submit(screen.getByTestId(/signin-form/i)));

        const errorEl = screen.getByText(/Invalid email or password./i);
        expect(errorEl).toBeInTheDocument();
    });

    it('should submit the form if the email and password are not empty', async () => {
        const emailInputEl = screen.getByPlaceholderText(/Email/i);
        expect(emailInputEl).toBeInTheDocument();

        const passwordInputEl = screen.getByPlaceholderText(/Password/i);
        expect(passwordInputEl).toBeInTheDocument();

        fireEvent.change(emailInputEl, {target: {value: 'test@example.com'}});
        fireEvent.change(passwordInputEl, {target: {value: 'Zaq12wsx'}});

        await act(() => fireEvent.submit(screen.getByTestId(/signin-form/i)))

        const errorEl = screen.queryByText(/Invalid email or password./i);
        expect(errorEl).not.toBeInTheDocument();
    });
});
