import React, { useState, useCallback, useId, useEffect } from "react";
import PropTypes from "prop-types";
import Form from "react-bootstrap/Form";
import FloatingLabel from "react-bootstrap/FloatingLabel";
import useFocus from "../hooks/useFocus";

const EditTodoForm = ({updateTodo, closeOnEsc, todo}) => {
    const [input, setInput] = useState(todo.taskText);
    const [inputRef, setInputFocus] = useFocus();
    const formID = useId();

    useEffect(() => setInputFocus(), []);

    const handleSubmit = useCallback((e) => {
        e.preventDefault();
        const text = input.trim();
        if (!text) {
            return false;
        }
        updateTodo(text, todo.id);
        setInput("");
    }, [updateTodo, todo, input]);

    return (
        <Form onSubmit={handleSubmit}>
            <FloatingLabel
                controlId={formID + "-todo-input"}
                label="Type your task here and press Enter"
            >
                <Form.Control
                    type='text'
                    ref={inputRef}
                    className='todo-input'
                    placeholder='Type your task here and press Enter'
                    onChange={(e) => setInput(e.target.value)}
                    onKeyUp={(e) => closeOnEsc(e)}
                    value={input}
                />
            </FloatingLabel>
        </Form>
    );
};

EditTodoForm.propTypes = {
    updateTodo: PropTypes.func.isRequired,
    closeOnEsc: PropTypes.func.isRequired,
    todo: PropTypes.object.isRequired
};

export default EditTodoForm;