
import React, { useRef, useState, useEffect, useCallback, useId } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import Alert from "react-bootstrap/Alert";
import axios from "../lib/axios";
import useInput from "../hooks/useInput";

const EMAIL_REGEX = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/;
const PASSWORD_REGEX = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.{8,})/;
const LOGIN_URL = "/login";

const SignInForm = () => {
    const formID = useId();
    const emailRef = useRef(null);

    const [validated, setValidated] = useState(false);

    const navigate = useNavigate();
    const location = useLocation();
    const from = location.state?.from?.pathname || "/";

    const [email, resetEmail, emailAttrs] = useInput("email", "");
    const [validEmail, setValidEmail] = useState(false);

    const [password, setPassword] = useState("");
    const [validPassword, setValidPassword] = useState(false);

    const [errMsg, setErrMsg] = useState("");

    useEffect(() => {
        const result = EMAIL_REGEX.test(email);
        setValidEmail(result);
    }, [email]);

    useEffect(() => {
        const result = PASSWORD_REGEX.test(password);
        setValidPassword(result);
    }, [password]);

    useEffect(() => {
        setErrMsg("");
    }, [email, password]);

    const handleSubmit = useCallback(async (e) => {
        e.preventDefault();

        const form = e.currentTarget;
        if (form.checkValidity() === false) {
            e.stopPropagation();
        }

        setValidated(true);

        if (EMAIL_REGEX.test(email) && PASSWORD_REGEX.test(password)) {
            try {
                await axios.post(
                    LOGIN_URL,
                    JSON.stringify({email, password}),
                    {
                        headers: { "Content-Type": "application/json" },
                        withCredentials: true
                    }
                );

                resetEmail("");
                setPassword("");
                navigate(from, {replace: true});
            } catch (error) {
                if (!error?.response) {
                    setErrMsg("Network error.");
                } else if (error.response?.status === 400) {
                    setErrMsg("Missing email or password.");
                } else if (error.response?.status === 401) {
                    setErrMsg("Invalid email or password.");
                } else {
                    setErrMsg("Something went wrong.");
                }
            }
        } else {
            setErrMsg("Invalid email or password.");
        }
    }, [email, password, from, resetEmail, navigate]);

    return (
        <>
            <Row>
                <Col md={{offset: 3, span: 6}}>
                    <h1>Sign In</h1>
                </Col>
            </Row>

            <Row className="mb-3">
                <Col md={{offset: 3, span: 6}}>
                    <Form onSubmit={handleSubmit} validated={validated}>
                        <Form.Group className="mb-3" controlId={formID + "-form-email"}>
                            <Form.Label>
                                Email address
                            </Form.Label>
                            <Form.Control
                                type="email"
                                ref={emailRef}
                                autoComplete="off"
                                placeholder="Email"
                                required
                                {...emailAttrs}
                                isValid={validEmail}
                                isInvalid={!validEmail && email.length > 0}
                            />
                        </Form.Group>

                        <Form.Group className="mb-3" controlId={formID + "-form-password"}>
                            <Form.Label>
                                Password
                            </Form.Label>
                            <Form.Control
                                type="password"
                                placeholder="Password"
                                value={password}
                                required
                                onChange={(e) => setPassword(e.target.value)}
                                isValid={validPassword}
                                isInvalid={!validPassword && password.length > 0}
                            />
                        </Form.Group>

                        { errMsg
                            ? <Alert variant="danger">{errMsg}</Alert>
                            : null
                        }

                        <Row>
                            <Col md={{offset: 3, span: 6}} className="d-grid">
                                <Button variant="primary" type="submit" disabled={!validEmail || !validPassword}>
                                    Sign In
                                </Button>
                            </Col>
                        </Row>
                    </Form>
                </Col>
            </Row>

            <Row className="mb-3">
                <Col md={{offset: 3, span: 6}} className="text-center">
                    <p>
                        Don&apos;t have an account yet?
                        {" "}
                        <Link to="/signup">Sign Up</Link>
                    </p>
                </Col>
            </Row>
        </>
    );
};

export default SignInForm;
