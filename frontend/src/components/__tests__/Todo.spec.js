import React from 'react';
import { render } from '@testing-library/react';
import Todo from '../Todo';

it('renders without crashing', () => {
    const todo = {
        id: 1,
        taskText: 'Learn React',
        isDone: false
    };
    const editTodo = jest.fn();
    const deleteTodo = jest.fn();
    const toggleComplete = jest.fn();

    const { getByTestId } = render(
        <Todo todo={todo} editTodo={editTodo} deleteTodo={deleteTodo} toggleComplete={toggleComplete}/>
    );
    expect(getByTestId(/todo-text-1/i)).toBeTruthy();
});
