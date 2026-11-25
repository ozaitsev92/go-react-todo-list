import { useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import Alert from "react-bootstrap/Alert";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import useAuth from "../auth/useAuth";
import { errorMessage } from "../lib/api";
import { safeRedirect } from "../lib/validation";

export default function SignInForm() {
    const { login } = useAuth();
    const navigate = useNavigate();
    const location = useLocation();

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [validated, setValidated] = useState(false);
    const [pending, setPending] = useState(false);
    const [error, setError] = useState("");

    const handleSubmit = async (e) => {
        e.preventDefault();
        setValidated(true);
        if (!e.currentTarget.checkValidity() || pending) {
            return;
        }
        setPending(true);
        setError("");
        try {
            await login(email.trim(), password);
            navigate(safeRedirect(location.state?.from?.pathname), { replace: true });
        } catch (err) {
            setError(errorMessage(err));
            setPending(false);
        }
    };

    return (
        <>
            <h1 className="h3 mb-4">Sign in</h1>
            <Form noValidate validated={validated} onSubmit={handleSubmit} aria-label="Sign in">
                <Form.Group className="mb-3" controlId="signin-email">
                    <Form.Label>Email</Form.Label>
                    <Form.Control
                        type="email"
                        autoComplete="username"
                        required
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                    />
                    <Form.Control.Feedback type="invalid">Enter your email.</Form.Control.Feedback>
                </Form.Group>
                <Form.Group className="mb-3" controlId="signin-password">
                    <Form.Label>Password</Form.Label>
                    <Form.Control
                        type="password"
                        autoComplete="current-password"
                        required
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                    <Form.Control.Feedback type="invalid">
                        Enter your password.
                    </Form.Control.Feedback>
                </Form.Group>
                {error && <Alert variant="danger">{error}</Alert>}
                <div className="d-grid">
                    <Button type="submit" disabled={pending}>
                        {pending ? "Signing in…" : "Sign in"}
                    </Button>
                </div>
            </Form>
            <p className="text-center mt-3">
                Don&apos;t have an account yet? <Link to="/signup">Sign up</Link>
            </p>
        </>
    );
}
