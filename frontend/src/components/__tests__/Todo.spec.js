import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import Todo from '../Todo';

describe('Todo', () => {
    const todo = {
        id: 1,
        taskText: 'Learn React',
        isDone: false
    };
    const editTodo = jest.fn();
    const deleteTodo = jest.fn();
    const toggleComplete = jest.fn();

    beforeEach(() => {
        render(<Todo todo={todo} editTodo={editTodo} deleteTodo={deleteTodo} toggleComplete={toggleComplete}/>);
    });

    it('should have a text field', () => {
        const todoEl = screen.getByText(/Learn React/i);
        expect(todoEl).toBeInTheDocument();
        fireEvent.click(todoEl);
        expect(toggleComplete).toHaveBeenCalledWith(1);
    });

    it('should have a edit button', () => {
        const editIconEl = screen.getByTestId(/todo-edit-1/i);
        expect(editIconEl).toBeInTheDocument();
        fireEvent.click(editIconEl);
        expect(editTodo).toHaveBeenCalledWith(1);
    });
    
    it('should have a delete button', () => {
        const deleteIconEl = screen.getByTestId(/todo-delete-1/i);
        expect(deleteIconEl).toBeInTheDocument();
        fireEvent.click(deleteIconEl);
        expect(deleteTodo).toHaveBeenCalledWith(1);
    });
});
