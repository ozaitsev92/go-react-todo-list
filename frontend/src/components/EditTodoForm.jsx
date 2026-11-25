import { useEffect, useId, useRef, useState } from "react";
import PropTypes from "prop-types";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import InputGroup from "react-bootstrap/InputGroup";
import { MAX_TASK_LENGTH } from "../lib/validation";
import { taskShape } from "./TodoItem";

// EditTodoForm edits a task's text: Enter saves, Escape or Cancel discards.
export default function EditTodoForm({ task, onSave, onCancel }) {
    const id = useId();
    const inputRef = useRef(null);
    const [text, setText] = useState(task.text);
    const [pending, setPending] = useState(false);

    useEffect(() => {
        inputRef.current?.focus();
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        const value = text.trim();
        if (!value || pending) {
            return;
        }
        setPending(true);
        await onSave(value);
        setPending(false);
    };

    return (
        <Form onSubmit={handleSubmit} aria-label={`Edit task "${task.text}"`}>
            <Form.Label htmlFor={id} visuallyHidden>
                Task text
            </Form.Label>
            <InputGroup>
                <Form.Control
                    id={id}
                    ref={inputRef}
                    type="text"
                    value={text}
                    maxLength={MAX_TASK_LENGTH}
                    onChange={(e) => setText(e.target.value)}
                    onKeyDown={(e) => e.key === "Escape" && onCancel()}
                    readOnly={pending}
                />
                <Button type="submit" disabled={pending || !text.trim()}>
                    Save
                </Button>
                <Button type="button" variant="outline-secondary" onClick={onCancel}>
                    Cancel
                </Button>
            </InputGroup>
        </Form>
    );
}

EditTodoForm.propTypes = {
    task: taskShape.isRequired,
    onSave: PropTypes.func.isRequired,
    onCancel: PropTypes.func.isRequired,
};
