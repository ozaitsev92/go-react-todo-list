import React, { useCallback, useState, useId } from "react";
import PropTypes from "prop-types";
import Form from "react-bootstrap/Form";
import FloatingLabel from "react-bootstrap/FloatingLabel";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";

const TodoForm = ({addTodo}) => {
    const [input, setInput] = useState("");
    const formID = useId();

    const handleSubmit = useCallback((e) => {
        e.preventDefault();
        const text = input.trim();
        if (!text) {
            return false;
        }
        addTodo({taskText: text});
        setInput("");
    }, [addTodo, input]);

    const clearOnEsc = useCallback((e) => {
        if (e.key === "Escape") {
            setInput("");
            e.target.blur();
        }
    });

    return (
        <Row>
            <Col md={{offset: 3, span: 6}}>
                <Form
                    onSubmit={handleSubmit}
                    data-testid="todo-form"
                >
                    <FloatingLabel
                        controlId={formID + "-todo-input"}
                        label="Type your task and press Enter or press ESC to cancel"
                    >
                        <Form.Control
                            type='text'
                            className='todo-input'
                            placeholder='Type your task and press Enter or press ESC to cancel'
                            onChange={(e) => setInput(e.target.value)}
                            onKeyUp={(e) => clearOnEsc(e)}
                            value={input}
                        />
                    </FloatingLabel>
                </Form>
            </Col>
        </Row>
    );
};

TodoForm.propTypes = {
    addTodo: PropTypes.func.isRequired
};

export default TodoForm;