import { useState } from "react";
import Alert from "react-bootstrap/Alert";
import Button from "react-bootstrap/Button";
import ListGroup from "react-bootstrap/ListGroup";
import Spinner from "react-bootstrap/Spinner";
import EditTodoForm from "../components/EditTodoForm";
import TodoForm from "../components/TodoForm";
import TodoItem from "../components/TodoItem";
import useTasks from "../hooks/useTasks";

export default function TodoList() {
    const {
        tasks,
        status,
        error,
        busyIds,
        reload,
        clearError,
        create,
        update,
        setCompleted,
        remove,
    } = useTasks();
    const [editingId, setEditingId] = useState(null);
    const [returnFocusId, setReturnFocusId] = useState(null);

    const stopEditing = (id) => {
        setEditingId(null);
        setReturnFocusId(id);
    };

    const save = async (id, text) => {
        if (await update(id, text)) {
            stopEditing(id);
        }
    };

    let body;
    if (status === "loading") {
        body = (
            <div className="text-center my-4">
                <Spinner animation="border" role="status">
                    <span className="visually-hidden">Loading tasks…</span>
                </Spinner>
            </div>
        );
    } else if (status === "error") {
        body = (
            <Button variant="outline-primary" onClick={reload}>
                Try again
            </Button>
        );
    } else if (tasks.length === 0) {
        body = <p className="text-muted text-center">No tasks yet. Add your first one above.</p>;
    } else {
        body = (
            <ListGroup as="ul" aria-label="Tasks">
                {tasks.map((task) => (
                    <ListGroup.Item as="li" key={task.id}>
                        {editingId === task.id ? (
                            <EditTodoForm
                                task={task}
                                onSave={(text) => save(task.id, text)}
                                onCancel={() => stopEditing(task.id)}
                            />
                        ) : (
                            <TodoItem
                                task={task}
                                busy={busyIds.has(task.id)}
                                focusEdit={returnFocusId === task.id}
                                onToggle={setCompleted}
                                onEdit={setEditingId}
                                onDelete={remove}
                            />
                        )}
                    </ListGroup.Item>
                ))}
            </ListGroup>
        );
    }

    return (
        <>
            <h1 className="h3 mb-4">What&apos;s the plan for today?</h1>
            {error && (
                <Alert variant="danger" dismissible onClose={clearError}>
                    {error}
                </Alert>
            )}
            <TodoForm onAdd={create} />
            {body}
        </>
    );
}
