import React from "react";
import Container from "react-bootstrap/Container";
import TodoWrapper from "../components/TodoWrapper";

const TodoList = () => {
    return (
        <Container className="mt-5">
            <TodoWrapper />
        </Container>
    );
};

export default TodoList;