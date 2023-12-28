import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import SignUpForm from '../SignUpForm';
import { BrowserRouter } from 'react-router-dom';
import { act } from 'react-dom/test-utils';

const MockedSignUpForm = () => {
    return (
        <BrowserRouter>
            <SignUpForm />
        </BrowserRouter>
    );
}

describe('SignUpForm', () => {
    beforeEach(() => {
        act(() => {
            render(<MockedSignUpForm />);
        });
    });

    it('should have a valid form', () => {
        const formEl = screen.getByTestId(/signup-form/i);
        expect(formEl).toBeInTheDocument();

        const emailLabelEl = screen.getByLabelText("Email");
        expect(emailLabelEl).toBeInTheDocument();

        const emailInputEl = screen.getByPlaceholderText("Email");
        expect(emailInputEl).toBeInTheDocument();

        const passwordLabelEl = screen.getByLabelText("Password");
        expect(passwordLabelEl).toBeInTheDocument();

        const passwordInputEl = screen.getByPlaceholderText("Password");
        expect(passwordInputEl).toBeInTheDocument();

        const confirmPasswordLabelEl = screen.getByLabelText("Confirm Password");
        expect(confirmPasswordLabelEl).toBeInTheDocument();

        const confirmPasswordInputEl = screen.getByPlaceholderText("Confirm Password");
        expect(confirmPasswordInputEl).toBeInTheDocument();

        const buttonEl = screen.getByRole("button", { name: /Sign Up/i });
        expect(buttonEl).toBeInTheDocument();
    });

    it('should not sign in if the email is empty', async () => {
        const emailInputEl = screen.getByPlaceholderText("Email");
        expect(emailInputEl).toBeInTheDocument();

        const passwordInputEl = screen.getByPlaceholderText("Password");
        expect(passwordInputEl).toBeInTheDocument();

        const confirmPasswordInputEl = screen.getByPlaceholderText("Confirm Password");
        expect(confirmPasswordInputEl).toBeInTheDocument();

        fireEvent.change(emailInputEl, {target: {value: ''}});
        fireEvent.change(passwordInputEl, {target: {value: 'Zaq12wsx'}});
        fireEvent.change(confirmPasswordInputEl, {target: {value: 'Zaq12wsx'}});

        await act(() => fireEvent.submit(screen.getByTestId(/signup-form/i)));

        const errorEl = screen.getByText(/Invalid email or password./i);
        expect(errorEl).toBeInTheDocument();
    });

    it('should not sign in if the password is empty', async () => {
        const emailInputEl = screen.getByPlaceholderText("Email");
        expect(emailInputEl).toBeInTheDocument();

        const passwordInputEl = screen.getByPlaceholderText("Password");
        expect(passwordInputEl).toBeInTheDocument();

        const confirmPasswordInputEl = screen.getByPlaceholderText("Confirm Password");
        expect(confirmPasswordInputEl).toBeInTheDocument();

        fireEvent.change(emailInputEl, {target: {value: 'test@example.com'}});
        fireEvent.change(passwordInputEl, {target: {value: ''}});
        fireEvent.change(confirmPasswordInputEl, {target: {value: ''}});

        await act(() => fireEvent.submit(screen.getByTestId(/signup-form/i)));

        const errorEl = screen.getByText(/Invalid email or password./i);
        expect(errorEl).toBeInTheDocument();
    });

    it('should submit the form if the email and password are not empty', async () => {
        const emailInputEl = screen.getByPlaceholderText("Email");
        expect(emailInputEl).toBeInTheDocument();

        const passwordInputEl = screen.getByPlaceholderText("Password");
        expect(passwordInputEl).toBeInTheDocument();

        const confirmPasswordInputEl = screen.getByPlaceholderText("Confirm Password");
        expect(confirmPasswordInputEl).toBeInTheDocument();

        fireEvent.change(emailInputEl, {target: {value: 'test@example.com'}});
        fireEvent.change(passwordInputEl, {target: {value: 'Zaq12wsx'}});
        fireEvent.change(confirmPasswordInputEl, {target: {value: 'Zaq12wsx'}});

        await act(() => fireEvent.submit(screen.getByTestId(/signup-form/i)))

        const errorEl = screen.queryByText(/Invalid email or password./i);
        expect(errorEl).not.toBeInTheDocument();
    });
});
