import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import TodoForm from '../TodoForm';

describe('TodoForm', () => {
    const addTodo = jest.fn();

    beforeEach(() => {
        render(<TodoForm addTodo={addTodo} />);
    });

    it('should have a valid form', () => {
        const formEl = screen.getByTestId(/todo-form/i);
        expect(formEl).toBeInTheDocument();

        const labelEl = screen.getByText(/Type your task and press Enter or press ESC to cancel/i);
        expect(labelEl).toBeInTheDocument();

        const inputEl = screen.getByPlaceholderText(/Type your task and press Enter or press ESC to cancel/i);
        expect(inputEl).toBeInTheDocument();
    });

    it('should not add a todo if the input is empty', () => {
        const inputEl = screen.getByPlaceholderText(/Type your task and press Enter or press ESC to cancel/i);
        expect(inputEl).toBeInTheDocument();

        inputEl.focus();
        fireEvent.change(inputEl, {target: {value: ''}})
        fireEvent.submit(screen.getByTestId(/todo-form/i));

        expect(addTodo).not.toHaveBeenCalled();
    });

    it('should add a todo if the input is not empty', () => {
        const inputEl = screen.getByPlaceholderText(/Type your task and press Enter or press ESC to cancel/i);
        expect(inputEl).toBeInTheDocument();

        inputEl.focus();
        fireEvent.change(inputEl, {target: {value: 'Learn React'}})
        fireEvent.submit(screen.getByTestId(/todo-form/i));

        expect(addTodo).toHaveBeenCalledWith({taskText: 'Learn React'});
    });
});
