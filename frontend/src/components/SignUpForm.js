
import React, { useRef, useState, useEffect } from "react";
import { faInfoCircle } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Link } from "react-router-dom";
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Alert from 'react-bootstrap/Alert';
import axios from "../api/axios";

const EMAIL_REGEX = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/;
const PASSWORD_REGEX = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.{8,})/;
const SIGNUP_URL = "/users";

const SignUpForm = () => {
    const emailRef = useRef(null);

    const [validated, setValidated] = useState(false);

    const [email, setEmail] = useState("");
    const [validEmail, setValidEmail] = useState(false);
    const [emailFocus, setEmailFocus] = useState(false);

    const [password, setPassword] = useState("");
    const [validPassword, setValidPassword] = useState(false);
    const [passwordFocus, setPasswordFocus] = useState(false);

    const [matchPassword, setMatchPassword] = useState("");
    const [validMatchPassword, setValidMatchPassword] = useState(false);
    const [matchPasswordFocus, setMatchPasswordFocus] = useState(false);

    const [errMsg, setErrMsg] = useState("");
    const [success, setSuccess] = useState(false);

    useEffect(() => {
        emailRef.current.focus();
    }, []);

    useEffect(() => {
        const result = EMAIL_REGEX.test(email);
        setValidEmail(result);
    }, [email]);

    useEffect(() => {
        const result = PASSWORD_REGEX.test(password);
        setValidPassword(result);
        const match = password === matchPassword;
        setValidMatchPassword(match);
    }, [password, matchPassword]);

    useEffect(() => {
        setErrMsg("");
    }, [email, password, matchPassword]);

    const handleSubmit = async (e) => {
        e.preventDefault();

        const form = e.currentTarget;
        if (form.checkValidity() === false) {
            e.stopPropagation();
        }

        setValidated(true);

        if (EMAIL_REGEX.test(email) && PASSWORD_REGEX.test(password) && password === matchPassword) {
            try {
                await axios.post(
                    SIGNUP_URL,
                    JSON.stringify({email, password}),
                    {
                        headers: { "Content-Type": "application/json" },
                        withCredentials: true
                    }
                );

                setEmail("");
                setPassword("");
                setMatchPassword("");
                setSuccess(true);
            } catch (error) {
                setSuccess(false);

                if (!error?.response) {
                    setErrMsg("Network error.");
                } else if (error.response?.status === 409) {
                    setErrMsg("Email already exists.");
                } else {
                    setErrMsg("Something went wrong.");
                }
            }
        } else {
            setErrMsg("Invalid email or password.");
            setSuccess(false);
        }
    };

    return (
        <>
            {
                success ? (
                    <>
                        <Row>
                            <Col md={{offset: 3, span: 6}}>
                                <h1>Success!</h1>
                            </Col>
                        </Row>
                        <Row>
                            <Col md={{offset: 3, span: 6}}>
                                <p>
                                    You have successfully signed up. Please <Link to="/signin">Sign In</Link>.
                                </p>
                            </Col>
                        </Row>
                    </>
                ): (
                    <>
                        <Row>
                            <Col md={{offset: 3, span: 6}}>
                                <h1>Sign Up</h1>
                            </Col>
                        </Row>
    
                        <Row className="mb-3">
                            <Col md={{offset: 3, span: 6}}>
                                <Form onSubmit={handleSubmit} validated={validated}>
                                    <Form.Group className="mb-3" controlId="email">
                                        <Form.Label>
                                            Email address
                                        </Form.Label>
                                        <Form.Control
                                            type="email"
                                            ref={emailRef}
                                            autoComplete="off"
                                            placeholder="Email"
                                            value={email}
                                            required
                                            onChange={(e) => setEmail(e.target.value)}
                                            onFocus={() => setEmailFocus(true)}
                                            onBlur={() => setEmailFocus(false)}
                                            isValid={validEmail}
                                            isInvalid={!validEmail && email.length > 0}
                                        />
                                        <Form.Text className="text-muted">
                                            {
                                                emailFocus && !validEmail && email.length > 0
                                                    ? <p className="info">
                                                        <FontAwesomeIcon icon={faInfoCircle} />
                                                        Email must be a valid email address.
                                                    </p>
                                                    : null
                                            }
                                        </Form.Text>
                                    </Form.Group>

                                    <Form.Group className="mb-3" controlId="formPassword">
                                        <Form.Label>
                                            Password
                                        </Form.Label>
                                        <Form.Control
                                            type="password"
                                            placeholder="Password"
                                            value={password}
                                            required
                                            onChange={(e) => setPassword(e.target.value)}
                                            onFocus={() => setPasswordFocus(true)}
                                            onBlur={() => setPasswordFocus(false)}
                                            isValid={validPassword}
                                            isInvalid={!validPassword && password.length > 0}
                                        />
                                        <Form.Text className="text-muted">
                                            {
                                                passwordFocus && !validPassword && password.length > 0
                                                    ? <p className="info">
                                                        <FontAwesomeIcon icon={faInfoCircle} />
                                                        Password must be at least 8 characters long, contain at least one lowercase letter, one uppercase letter, and one number.
                                                    </p>
                                                    : null
                                            }
                                        </Form.Text>
                                    </Form.Group>

                                    <Form.Group className="mb-3" controlId="formMatchPassword">
                                        <Form.Label>
                                            Confirm Password
                                        </Form.Label>
                                        <Form.Control
                                            type="password"
                                            placeholder="Confirm Password"
                                            value={matchPassword}
                                            required
                                            onChange={(e) => setMatchPassword(e.target.value)}
                                            onFocus={() => setMatchPasswordFocus(true)}
                                            onBlur={() => setMatchPasswordFocus(false)}
                                            isValid={validMatchPassword && matchPassword}
                                            isInvalid={!validMatchPassword && !matchPassword.length}
                                        />
                                        <Form.Text className="text-muted">
                                            {
                                                matchPasswordFocus && !validMatchPassword && matchPassword.length > 0
                                                    ? <p className="info">
                                                        <FontAwesomeIcon icon={faInfoCircle} />
                                                        Passwords must match.
                                                    </p>
                                                    : null
                                            }
                                        </Form.Text>
                                    </Form.Group>

                                    {
                                        errMsg
                                            ? <Alert variant="danger">{errMsg}</Alert>
                                            : null
                                    }

                                    <Row>
                                        <Col md={{offset: 3, span: 6}} className="d-grid">
                                            <Button variant="primary" type="submit" disabled={!validEmail || !validPassword || !validMatchPassword}>
                                                Sign Up
                                            </Button>
                                        </Col>
                                    </Row>
                                </Form>
                            </Col>
                        </Row>

                        <Row className="mb-3">
                            <Col md={{offset: 3, span: 6}} className="text-center">
                                <p>
                                    Already have an account?
                                    {" "}
                                    <Link to="/signin">Sign In</Link>
                                </p>
                            </Col>
                        </Row>
                    </>
                )
            }
        </>
    );
};

export default SignUpForm;
