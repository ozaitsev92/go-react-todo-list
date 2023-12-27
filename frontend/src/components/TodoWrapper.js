import React, {useState, useEffect, useCallback} from "react";
import TodoForm from "./TodoForm";
import Todo from "./Todo";
import EditTodoForm from "./EditTodoForm";
import axios from "../lib/axios";
import useAuth from "../hooks/useAuth";
import { useNavigate } from "react-router-dom";
import Alert from "react-bootstrap/Alert";
import Button from "react-bootstrap/Button";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import ListGroup from "react-bootstrap/ListGroup";

const TodoWrapper = () => {
    const [errMsg, setErrMsg] = useState("");
    const navigate = useNavigate();
    const { auth } = useAuth();
    const [todos, setTodos] = useState([]);

    useEffect(() => {

        let isMounted = true;
        const controller = new AbortController();

        const getTodos = async () => {
            setErrMsg("");
            const userID = auth?.user?.id;

            try {
                const response = await axios.get(`/users/${userID}/tasks`, {
                    headers: { "Content-Type": "application/json" },
                    withCredentials: true,
                    signal: controller.signal
                });

                if (isMounted) {
                    setTodos(response.data || []);
                }
            } catch (error) {
                if (error?.response) {
                    if (error.response?.status === 401 || error.response?.status === 403) {
                        navigate("/signin");
                    } else {
                        setErrMsg("Something went wrong.");
                    }
                }
            }
        };

        getTodos();

        return () => {
            controller.abort();
            isMounted = false;
        };
    }, [auth, navigate]);

    const addTodo = useCallback(async (todo) => {
        setErrMsg("");

        const userID = auth?.user?.id;

        let taskOrder = -1;
        for (let i = 0; i < todos.length; i++) {
            if (todos[i].taskOrder > taskOrder) {
                taskOrder = todos[i].taskOrder;
            }
        }
        taskOrder += 1;

        const newTodo = {
            taskText: todo.taskText,
            userId: userID,
            taskOrder
        };

        try {
            await axios.post(
                `/users/${userID}/tasks`,
                JSON.stringify(newTodo),
                {
                    headers: { "Content-Type": "application/json" },
                    withCredentials: true
                }
            );

            const response = await axios.get(`/users/${userID}/tasks`, {
                headers: { "Content-Type": "application/json" },
                withCredentials: true
            });

            setTodos(response.data || []);
        } catch (error) {
            if (error?.response) {
                if (error.response?.status === 401 || error.response?.status === 403) {
                    navigate("/signin");
                } else {
                    setErrMsg("Something went wrong.");
                }
            }
        }
    }, [auth, todos, navigate]);

    const toggleComplete = useCallback(async (id) => {
        setErrMsg("");

        const userID = auth?.user?.id;
        const todo = todos.slice()
            .filter((todo) => todo.id === id)[0];

        if (todo) {
            try {
                const url = todo.isDone
                    ? `/users/${userID}/tasks/${id}/mark-not-done`
                    : `/users/${userID}/tasks/${id}/mark-done`;

                await axios.put(
                    url,
                    null,
                    {
                        headers: { "Content-Type": "application/json" },
                        withCredentials: true
                    }
                );

                const response = await axios.get(`/users/${userID}/tasks`, {
                    headers: { "Content-Type": "application/json" },
                    withCredentials: true
                });

                setTodos(response.data || []);
            } catch (error) {
                if (error?.response) {
                    if (error.response?.status === 401 || error.response?.status === 403) {
                        navigate("/signin");
                    } else {
                        setErrMsg("Something went wrong.");
                    }
                }
            }
        }
    }, [auth, todos, navigate]);

    const deleteTodo = useCallback(async (id) => {
        setErrMsg("");

        const userID = auth?.user?.id;

        try {
            await axios.delete(
                `/users/${userID}/tasks/${id}`,
                {
                    headers: { "Content-Type": "application/json" },
                    withCredentials: true
                }
            );

            const response = await axios.get(`/users/${userID}/tasks`, {
                headers: { "Content-Type": "application/json" },
                withCredentials: true
            });

            setTodos(response.data || []);
        } catch (error) {
            if (error?.response) {
                if (error.response?.status === 401 || error.response?.status === 403) {
                    navigate("/signin");
                } else {
                    setErrMsg("Something went wrong.");
                }
            }
        }
    }, [auth, navigate]);

    const editTodo = useCallback((id) => {
        const updatedTodos = todos.slice()
            .map((todo) => {
                if (todo.id === id) {
                    todo.isEditing = !todo.isEditing;
                } else {
                    todo.isEditing = false;
                }
                return todo;
            });
        setTodos(updatedTodos);
    }, [todos]);

    const closeOnEsc = useCallback((e) => {
        if (e.key === "Escape") {
            const updatedTodos = todos.slice()
                .map((todo) => {
                    todo.isEditing = false;
                    return todo;
                });
            setTodos(updatedTodos);
        }
    }, [todos]);

    const updateTodo = useCallback(async (input, id) => {
        setErrMsg("");

        const userID = auth?.user?.id;
        const todo = todos.slice()
            .filter((todo) => todo.id === id)[0];

        if (todo) {
            try {
                todo.taskText = input;

                await axios.put(
                    `/users/${userID}/tasks/${id}`,
                    JSON.stringify(todo),
                    {
                        headers: { "Content-Type": "application/json" },
                        withCredentials: true
                    }
                );

                const response = await axios.get(`/users/${userID}/tasks`, {
                    headers: { "Content-Type": "application/json" },
                    withCredentials: true
                });

                setTodos(response.data || []);
            } catch (error) {
                if (error?.response) {
                    if (error.response?.status === 401 || error.response?.status === 403) {
                        navigate("/signin");
                    } else {
                        setErrMsg("Something went wrong.");
                    }
                }
            }
        }
    }, [auth, todos, navigate]);

    const logout = useCallback(async () => {
        setErrMsg("");

        try {
            await axios.post(
                "/logout",
                null,
                {
                    headers: { "Content-Type": "application/json" },
                    withCredentials: true
                }
            );

            navigate("/signin");
        } catch (error) {
            if (error.response?.status === 401 || error.response?.status === 403) {
                navigate("/signin");
            } else {
                setErrMsg("Something went wrong.");
            }
        }
    }, [navigate]);

    return (
        <>
            <Row>
                <Col md={{offset: 3, span: 6}} className="text-center">
                    <h1>What&apos;s the Plan for Today?</h1>
                </Col>
            </Row>

            <Row>
                <Col md={{offset: 3, span: 6}}>
                    {
                        errMsg
                            ? <Alert variant="danger">{errMsg}</Alert>
                            : null
                    }
                </Col>
            </Row>

            <TodoForm addTodo={addTodo} />

            <Row>
                <Col md={{offset: 3, span: 6}}>
                    <hr />
                </Col>
            </Row>

            <Row className="mb-3">
                <Col md={{offset: 3, span: 6}}>
                    <ListGroup>
                        {todos.map((todo) => (
                            todo.isEditing
                                ? (
                                    <ListGroup.Item
                                        key={todo.id}
                                        data-testid={`todo-item-${todo.id}`}
                                    >
                                        <EditTodoForm todo={todo} key={todo.id} updateTodo={updateTodo} closeOnEsc={closeOnEsc} />
                                    </ListGroup.Item>
                                )
                                : (
                                    <ListGroup.Item
                                        key={todo.id}
                                        data-testid={`todo-item-${todo.id}`}
                                    >
                                        <Todo todo={todo} toggleComplete={toggleComplete} deleteTodo={deleteTodo} editTodo={editTodo} />
                                    </ListGroup.Item>
                                )
                        ))}
                    </ListGroup>
                </Col>
            </Row>

            <Row className="mb-3">
                <Col md={{offset: 3, span: 6}} className="text-center">
                    <Button variant="link" className="cursor-pointer" onClick={logout}>logout</Button>
                </Col>
            </Row>
        </>
    );
};

export default TodoWrapper;