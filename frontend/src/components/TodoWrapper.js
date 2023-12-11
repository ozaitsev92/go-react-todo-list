import React, {useState, useEffect} from "react";
import TodoForm from "./TodoForm";
import Todo from "./Todo";
import EditTodoForm from "./EditTodoForm";
import axios from "../api/axios";
import useAuth from "../hooks/useAuth";
import { useNavigate } from "react-router-dom";

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

    const addTodo = async (todo) => {
        setErrMsg("");

        const userID = auth?.user?.id;

        const newTodo = {
            task_text: todo.task_text,
            user_id: userID,
            task_order: 0 //todo: implement task_order
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
    };

    const toggleComplete = async (id) => {
        setErrMsg("");

        const userID = auth?.user?.id;
        const todo = todos.clice()
            .filter((todo) => todo.id === id)[0];

        if (todo) {
            try {
                const url = todo.is_done
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
    };

    const deleteTodo = async (id) => {
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
    };

    const editTodo = (id) => {
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
    };

    const updateTodo = async (input, id) => {
        setErrMsg("");

        const userID = auth?.user?.id;
        const todo = todos.slice()
            .filter((todo) => todo.id === id)[0];

        if (todo) {
            try {
                todo.task_text = input;

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
    };

    const logout = async () => {
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
    };

    return (
        <div className='TodoWrapper'>
            <h1>What&apos;s the Plan for Today?</h1>
            {
                errMsg
                    ? <p className="error">{errMsg}</p>
                    : null
            }
            <TodoForm addTodo={addTodo} />
            {todos.map((todo) => (
                todo.isEditing
                    ? <EditTodoForm todo={todo} key={todo.id} updateTodo={updateTodo} />
                    : <Todo todo={todo} key={todo.id} toggleComplete={toggleComplete} deleteTodo={deleteTodo} editTodo={editTodo} />
            ))}
            <button type="link" onClick={logout}>logout</button>
        </div>
    );
};

export default TodoWrapper;