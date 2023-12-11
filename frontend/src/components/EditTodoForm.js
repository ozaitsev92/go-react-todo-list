import React, {useState} from "react";
import PropTypes from "prop-types";

const EditTodoForm = ({updateTodo, todo}) => {
    const [input, setInput] = useState(todo.task_text);

    const handleSubmit = (e) => {
        e.preventDefault();
        const text = input.trim();
        if (!text) {
            return false;
        }
        updateTodo(text, todo.id);
        setInput("");
    };

    return (
        <form
            className='TodoForm'
            onSubmit={handleSubmit}
        >
            <input
                type='text'
                className='todo-input'
                placeholder='Update task'
                onChange={(e) => setInput(e.target.value)}
                value={input}
            />
            <button
                type='submit'
                className='todo-btn'
            >
                Update Task
            </button>
        </form>
    );
};

EditTodoForm.propTypes = {
    updateTodo: PropTypes.func.isRequired,
    todo: PropTypes.object.isRequired
};

export default EditTodoForm;