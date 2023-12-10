
import React, { useRef, useState, useEffect } from "react";
import { faCheck, faTimes, faInfoCircle } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "../api/axios";

const EMAIL_REGEX = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/;
const PASSWORD_REGEX = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.{8,})/;

const SignUpForm = () => {
    const emailRef = useRef(null);

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

        if (EMAIL_REGEX.test(email) && PASSWORD_REGEX.test(password) && password === matchPassword) {
            //todo: send request to backend
            try {
                const response = await axios.post(
                    "/users",
                    JSON.stringify({email, password}),
                    {
                        headers: { "Content-Type": "application/json" },
                        withCredentials: true
                    }
                );
                console.log(JSON.stringify(response));
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
                    setErrMsg("An error occurred.");
                }
            }
        } else {
            setErrMsg("Invalid email or password.");
        }
    };

    return (
        <>
            {
                success ? (
                    <section>
                        <h1>Success!</h1>
                        <p>
                            You have successfully signed up. Please <a href="/signin">Sign In</a>.
                        </p>
                    </section>
                ): (
                    <section>
                        <h1>Sign Up</h1>
                        {
                            errMsg
                                ? <p className="error">{errMsg}</p>
                                : null
                        }
                        <form onSubmit={handleSubmit}>
                            <div>
                                <label htmlFor="email">
                                    Email:
                                    <span className={validEmail ? "valid" : "hide"}>
                                        <FontAwesomeIcon icon={faCheck} />
                                    </span>
                                    <span className={validEmail || !email ? "hide" : "invalid"}>
                                        <FontAwesomeIcon icon={faTimes} />
                                    </span>
                                </label>
                                <input
                                    type="email"
                                    id="email"
                                    ref={emailRef}
                                    autoComplete="off"
                                    placeholder="Email"
                                    value={email}
                                    required
                                    onChange={(e) => setEmail(e.target.value)}
                                    onFocus={() => setEmailFocus(true)}
                                    onBlur={() => setEmailFocus(false)}
                                />
                            </div>
                            {
                                emailFocus && !validEmail && email.length > 0
                                    ? <p className="info">
                                        <FontAwesomeIcon icon={faInfoCircle} />
                                        Email must be a valid email address.
                                    </p>
                                    : null
                            }
                            <div>
                                <label htmlFor="password">
                                    Password:
                                    <span className={validPassword ? "valid" : "hide"}>
                                        <FontAwesomeIcon icon={faCheck} />
                                    </span>
                                    <span className={validPassword || !password ? "hide" : "invalid"}>
                                        <FontAwesomeIcon icon={faTimes} />
                                    </span>
                                </label>
                                <input
                                    type="password"
                                    id="password"
                                    placeholder="Password"
                                    value={password}
                                    required
                                    onChange={(e) => setPassword(e.target.value)}
                                    onFocus={() => setPasswordFocus(true)}
                                    onBlur={() => setPasswordFocus(false)}
                                />
                            </div>
                            {
                                passwordFocus && !validPassword && password.length > 0
                                    ? <p className="info">
                                        <FontAwesomeIcon icon={faInfoCircle} />
                                        Password must be at least 8 characters long, contain at least one lowercase letter, one uppercase letter, and one number.
                                    </p>
                                    : null
                            }
                            <div>
                                <label htmlFor="match-password">
                                    Confirm Password:
                                    <span className={validMatchPassword && matchPassword ? "valid" : "hide"}>
                                        <FontAwesomeIcon icon={faCheck} />
                                    </span>
                                    <span className={!validMatchPassword || !matchPassword ? "invalid" : "hide"}>
                                        <FontAwesomeIcon icon={faTimes} />
                                    </span>
                                </label>
                                <input
                                    type="password"
                                    id="match-password"
                                    placeholder="Confirm Password"
                                    value={matchPassword}
                                    required
                                    onChange={(e) => setMatchPassword(e.target.value)}
                                    onFocus={() => setMatchPasswordFocus(true)}
                                    onBlur={() => setMatchPasswordFocus(false)}
                                />
                            </div>
                            {
                                matchPasswordFocus && !validMatchPassword && matchPassword.length > 0
                                    ? <p className="info">
                                        <FontAwesomeIcon icon={faInfoCircle} />
                                        Passwords must match.
                                    </p>
                                    : null
                            }
                            <button
                                disabled={!validEmail || !validPassword || !validMatchPassword}
                                type="submit"
                            >
                                Sign Up
                            </button>
                        </form>
                        <p>
                            Already have an account? <a href="/signin">Sign In</a>
                        </p>
                    </section>
                )
            }
        </>
    );
};

export default SignUpForm;
