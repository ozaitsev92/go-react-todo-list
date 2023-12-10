import React, {useState, useEffect} from "react";
import TodoForm from "./TodoForm";
import Todo from "./Todo";
import EditTodoForm from "./EditTodoForm";
import axios from "../api/axios";

const TodoWrapper = () => {
    const [todos, setTodos] = useState([]);

    useEffect(() => {
        let isMounted = true;
        const controller = new AbortController();

        const getTodos = async () => {
            const userID = "1234567890";
            try {
                const response = await axios.get(`/users/${userID}/tasks`, {
                    signal: controller.signal
                });
                if (isMounted) {
                    setTodos(response.data);
                }
            } catch (error) {
                console.log(error);
            }
        };

        getTodos();

        return () => {
            controller.abort();
            isMounted = false;
        };
    }, []);

    const addTodo = (todo) => {
        const newTodos = [{
            id: todos.length + 1,
            text: todo.text,
            completed: false,
            isEditing: false
        }, ...todos];
        setTodos(newTodos);
    };

    const toggleComplete = (id) => {
        const updatedTodos = todos.slice()
            .map((todo) => {
                if (todo.id === id) {
                    todo.completed = !todo.completed;
                }
                return todo;
            });
        setTodos(updatedTodos);
    };

    const deleteTodo = (id) => {
        const updatedTodos = todos.filter((todo) => todo.id !== id);
        setTodos(updatedTodos);
    };

    const editTodo = (id) => {
        const updatedTodos = todos.slice()
            .map((todo) => {
                if (todo.id === id) {
                    todo.isEditing = !todo.isEditing;
                }
                return todo;
            });
        setTodos(updatedTodos);
    };

    const updateTodo = (input, id) => {
        const updatedTodos = todos.slice()
            .map((todo) => {
                if (todo.id === id) {
                    todo.text = input;
                    todo.isEditing = !todo.isEditing;
                }
                return todo;
            });
        setTodos(updatedTodos);
    };

    return (
        <div className='TodoWrapper'>
            <h1>What&apos;s the Plan for Today?</h1>
            <TodoForm addTodo={addTodo} />
            {todos.map((todo) => (
                todo.isEditing
                    ? <EditTodoForm todo={todo} key={todo.id} updateTodo={updateTodo} />
                    : <Todo todo={todo} key={todo.id} toggleComplete={toggleComplete} deleteTodo={deleteTodo} editTodo={editTodo} />
            ))}
        </div>
    );
};

export default TodoWrapper;