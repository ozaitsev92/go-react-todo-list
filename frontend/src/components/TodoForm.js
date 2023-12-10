import React, {useState} from "react";
import PropTypes from "prop-types";

const TodoForm = ({addTodo}) => {
    const [input, setInput] = useState("");

    const handleSubmit = (e) => {
        e.preventDefault();
        const text = input.trim();
        if (!text) {
            return false;
        }
        addTodo({text});
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
                placeholder='What is the task today?'
                onChange={(e) => setInput(e.target.value)}
                value={input}
            />
            <button
                type='submit'
                className='todo-btn'
            >
                Add Task
            </button>
        </form>
    );
};

TodoForm.propTypes = {
    addTodo: PropTypes.func.isRequired
};

export default TodoForm;