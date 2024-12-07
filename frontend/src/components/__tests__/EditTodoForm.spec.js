import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import EditTodoForm from '../EditTodoForm';

describe('EditTodoForm', () => {
    const todo = {
        id: 1,
        text: 'Lorem Ipsum',
        completed: false,
        isEditing: true
    };
    const updateTodo = jest.fn();
    const closeOnEsc = jest.fn();

    beforeEach(() => {
        render(<EditTodoForm todo={todo} updateTodo={updateTodo} closeOnEsc={closeOnEsc} />);
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
        fireEvent.change(inputEl, { target: { value: '' } })
        fireEvent.submit(screen.getByTestId(/todo-form/i));

        expect(updateTodo).not.toHaveBeenCalled();
    });

    it('should add a todo if the input is not empty', () => {
        const inputEl = screen.getByPlaceholderText(/Type your task and press Enter or press ESC to cancel/i);
        expect(inputEl).toBeInTheDocument();

        inputEl.focus();
        fireEvent.change(inputEl, { target: { value: 'Learn React' } })
        fireEvent.submit(screen.getByTestId(/todo-form/i));

        expect(updateTodo).toHaveBeenCalledWith('Learn React', todo.id);
    });

    it('should close the form if the ESC key is pressed', () => {
        const inputEl = screen.getByPlaceholderText(/Type your task and press Enter or press ESC to cancel/i);
        expect(inputEl).toBeInTheDocument();

        inputEl.focus();
        fireEvent.keyUp(inputEl, { key: 'Escape', code: 'Escape', keyCode: 27, charCode: 27 });

        expect(closeOnEsc).toHaveBeenCalled();
    });
});
