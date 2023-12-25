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

    return (
        <Row>
            <Col md={{offset: 3, span: 6}}>
                <Form onSubmit={handleSubmit}>
                    <FloatingLabel
                        controlId={formID + "-todo-input"}
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
            </Col>
        </Row>
    );
};

TodoForm.propTypes = {
    addTodo: PropTypes.func.isRequired
};

export default TodoForm;