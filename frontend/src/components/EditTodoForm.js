import React, {useState} from "react";
import PropTypes from "prop-types";
import Form from 'react-bootstrap/Form';
import FloatingLabel from 'react-bootstrap/FloatingLabel';

const EditTodoForm = ({updateTodo, todo}) => {
    const [input, setInput] = useState(todo.taskText);

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
        <Form onSubmit={handleSubmit}>
            <FloatingLabel
                controlId="todo-input"
                label="Type your task here and press Enter"
            >
                <Form.Control
                    type='text'
                    className='todo-input'
                    placeholder='Type your task here and press Enter'
                    onChange={(e) => setInput(e.target.value)}
                    value={input}
                />
            </FloatingLabel>
        </Form>
    );
};

EditTodoForm.propTypes = {
    updateTodo: PropTypes.func.isRequired,
    todo: PropTypes.object.isRequired
};

export default EditTodoForm;