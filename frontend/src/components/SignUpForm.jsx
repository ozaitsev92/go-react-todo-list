import { useState } from "react";
import { Link } from "react-router-dom";
import Alert from "react-bootstrap/Alert";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import { authApi, errorMessage } from "../lib/api";
import {
    MAX_PASSWORD_BYTES,
    MIN_PASSWORD_BYTES,
    isValidEmail,
    isValidPassword,
} from "../lib/validation";

export default function SignUpForm() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirm, setConfirm] = useState("");
    const [submitted, setSubmitted] = useState(false);
    const [pending, setPending] = useState(false);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState(false);

    const emailOk = isValidEmail(email);
    const passwordOk = isValidPassword(password);
    const confirmOk = confirm === password && confirm !== "";

    // Flag a field once it has input or a submit was attempted.
    const invalid = (ok, value) => !ok && (submitted || value !== "");

    const handleSubmit = async (e) => {
        e.preventDefault();
        setSubmitted(true);
        if (!emailOk || !passwordOk || !confirmOk || pending) {
            return;
        }
        setPending(true);
        setError("");
        try {
            await authApi.register(email.trim(), password);
            setSuccess(true);
        } catch (err) {
            setError(errorMessage(err));
        } finally {
            setPending(false);
        }
    };

    if (success) {
        return (
            <>
                <h1 className="h3 mb-4">Account created</h1>
                <p role="status">
                    You can now <Link to="/signin">sign in</Link>.
                </p>
            </>
        );
    }

    return (
        <>
            <h1 className="h3 mb-4">Sign up</h1>
            <Form noValidate onSubmit={handleSubmit} aria-label="Sign up">
                <Form.Group className="mb-3" controlId="signup-email">
                    <Form.Label>Email</Form.Label>
                    <Form.Control
                        type="email"
                        autoComplete="username"
                        required
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        isInvalid={invalid(emailOk, email)}
                        isValid={emailOk}
                    />
                    <Form.Control.Feedback type="invalid">
                        Enter a valid email address.
                    </Form.Control.Feedback>
                </Form.Group>
                <Form.Group className="mb-3" controlId="signup-password">
                    <Form.Label>Password</Form.Label>
                    <Form.Control
                        type="password"
                        autoComplete="new-password"
                        required
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        isInvalid={invalid(passwordOk, password)}
                        isValid={passwordOk}
                        aria-describedby="signup-password-help"
                    />
                    <Form.Text id="signup-password-help" muted>
                        {MIN_PASSWORD_BYTES} to {MAX_PASSWORD_BYTES} characters.
                    </Form.Text>
                </Form.Group>
                <Form.Group className="mb-3" controlId="signup-confirm">
                    <Form.Label>Confirm password</Form.Label>
                    <Form.Control
                        type="password"
                        autoComplete="new-password"
                        required
                        value={confirm}
                        onChange={(e) => setConfirm(e.target.value)}
                        isInvalid={invalid(confirmOk, confirm)}
                        isValid={confirmOk}
                    />
                    <Form.Control.Feedback type="invalid">
                        Passwords must match.
                    </Form.Control.Feedback>
                </Form.Group>
                {error && <Alert variant="danger">{error}</Alert>}
                <div className="d-grid">
                    <Button type="submit" disabled={pending}>
                        {pending ? "Creating account…" : "Sign up"}
                    </Button>
                </div>
            </Form>
            <p className="text-center mt-3">
                Already have an account? <Link to="/signin">Sign in</Link>
            </p>
        </>
    );
}
