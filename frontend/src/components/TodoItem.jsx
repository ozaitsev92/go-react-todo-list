import { useEffect, useRef } from "react";
import PropTypes from "prop-types";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPenToSquare, faTrashCan } from "@fortawesome/free-regular-svg-icons";

export const taskShape = PropTypes.shape({
    id: PropTypes.string.isRequired,
    text: PropTypes.string.isRequired,
    completed: PropTypes.bool.isRequired,
});

// TodoItem shows one task with native, keyboard-operable controls.
export default function TodoItem({ task, busy, focusEdit, onToggle, onEdit, onDelete }) {
    const editRef = useRef(null);

    // Return focus to the Edit button after leaving edit mode.
    useEffect(() => {
        if (focusEdit) {
            editRef.current?.focus();
        }
    }, [focusEdit]);

    return (
        <div className="d-flex align-items-center gap-2">
            <Form.Check id={`task-${task.id}`} className="flex-grow-1 mb-0">
                <Form.Check.Input
                    type="checkbox"
                    checked={task.completed}
                    disabled={busy}
                    onChange={() => onToggle(task.id, !task.completed)}
                />
                <Form.Check.Label className={task.completed ? "completed" : undefined}>
                    {task.text}
                </Form.Check.Label>
            </Form.Check>
            <Button
                ref={editRef}
                variant="link"
                size="sm"
                aria-label={`Edit "${task.text}"`}
                disabled={busy}
                onClick={() => onEdit(task.id)}
            >
                <FontAwesomeIcon icon={faPenToSquare} />
            </Button>
            <Button
                variant="link"
                size="sm"
                className="text-danger"
                aria-label={`Delete "${task.text}"`}
                disabled={busy}
                onClick={() => onDelete(task.id)}
            >
                <FontAwesomeIcon icon={faTrashCan} />
            </Button>
        </div>
    );
}

TodoItem.propTypes = {
    task: taskShape.isRequired,
    busy: PropTypes.bool,
    focusEdit: PropTypes.bool,
    onToggle: PropTypes.func.isRequired,
    onEdit: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
};
