import React from "react";
import PropTypes from "prop-types";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPenToSquare } from "@fortawesome/free-regular-svg-icons";
import { faTrashAlt } from "@fortawesome/free-regular-svg-icons";

const Todo = ({todo, toggleComplete, editTodo, deleteTodo}) => {
    return (
        <div className='Todo'>
            <p
                className={`${todo.is_done ? "completed" : ""}`}
                onClick={() => toggleComplete(todo.id)}
            >
                {todo.task_text}
            </p>
            <div>
                <FontAwesomeIcon icon={faPenToSquare} onClick={() => editTodo(todo.id)} />
                <FontAwesomeIcon icon={faTrashAlt} onClick={() => deleteTodo(todo.id)} />
            </div>
        </div>
    );
};

Todo.propTypes = {
    todo: PropTypes.object.isRequired,
    toggleComplete: PropTypes.func.isRequired,
    editTodo: PropTypes.func.isRequired,
    deleteTodo: PropTypes.func.isRequired
};

export default Todo;