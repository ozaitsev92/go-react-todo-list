import { useId, useState } from "react";
import PropTypes from "prop-types";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import InputGroup from "react-bootstrap/InputGroup";
import { MAX_TASK_LENGTH } from "../lib/validation";

// TodoForm adds a task, clearing the input only after onAdd succeeds.
export default function TodoForm({ onAdd }) {
    const id = useId();
    const [text, setText] = useState("");
    const [pending, setPending] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        const value = text.trim();
        if (!value || pending) {
            return;
        }
        setPending(true);
        const ok = await onAdd(value);
        setPending(false);
        if (ok) {
            setText("");
        }
    };

    return (
        <Form onSubmit={handleSubmit} aria-label="Add task" className="mb-4">
            <Form.Label htmlFor={id} visuallyHidden>
                New task
            </Form.Label>
            <InputGroup>
                <Form.Control
                    id={id}
                    type="text"
                    placeholder="What needs to be done?"
                    value={text}
                    maxLength={MAX_TASK_LENGTH}
                    onChange={(e) => setText(e.target.value)}
                    onKeyDown={(e) => e.key === "Escape" && setText("")}
                    readOnly={pending}
                />
                <Button type="submit" disabled={pending || !text.trim()}>
                    {pending ? "Adding…" : "Add"}
                </Button>
            </InputGroup>
        </Form>
    );
}

TodoForm.propTypes = {
    onAdd: PropTypes.func.isRequired,
};
