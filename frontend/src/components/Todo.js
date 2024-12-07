import React from "react";
import PropTypes from "prop-types";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPenToSquare } from "@fortawesome/free-regular-svg-icons";
import { faTrashAlt } from "@fortawesome/free-regular-svg-icons";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";

const Todo = ({ todo, toggleComplete, editTodo, deleteTodo }) => {
    return (
        <Row style={{ justifyContent: "space-between" }}>
            <Col sm="auto">
                <span
                    className={`${todo.completed ? "completed" : ""} cursor-pointer`}
                    onClick={() => toggleComplete(todo.id)}
                    data-testid={`todo-text-${todo.id}`}
                >
                    {todo.text}
                </span>
            </Col>
            <Col sm="auto">
                <span>
                    <FontAwesomeIcon
                        className="cursor-pointer"
                        icon={faPenToSquare}
                        onClick={() => editTodo(todo.id)}
                        data-testid={`todo-edit-${todo.id}`}
                    />
                    {" "}
                    <FontAwesomeIcon
                        className="cursor-pointer"
                        icon={faTrashAlt}
                        onClick={() => deleteTodo(todo.id)}
                        data-testid={`todo-delete-${todo.id}`}
                    />
                </span>
            </Col>
        </Row>
    );
};

Todo.propTypes = {
    todo: PropTypes.object.isRequired,
    toggleComplete: PropTypes.func.isRequired,
    editTodo: PropTypes.func.isRequired,
    deleteTodo: PropTypes.func.isRequired
};

export default Todo;