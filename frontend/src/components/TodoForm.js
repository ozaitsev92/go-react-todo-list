import React, {useState} from "react";
import PropTypes from "prop-types";
import Form from 'react-bootstrap/Form';
import FloatingLabel from 'react-bootstrap/FloatingLabel';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';

const TodoForm = ({addTodo}) => {
    const [input, setInput] = useState("");

    const handleSubmit = (e) => {
        e.preventDefault();
        const text = input.trim();
        if (!text) {
            return false;
        }
        addTodo({taskText: text});
        setInput("");
    };

    return (
        <Row>
            <Col md={{offset: 3, span: 6}}>
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
            </Col>
        </Row>
    );
};

TodoForm.propTypes = {
    addTodo: PropTypes.func.isRequired
};

export default TodoForm;