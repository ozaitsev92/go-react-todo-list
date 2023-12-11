
import React, { useRef, useState, useEffect } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { faCheck, faTimes, faInfoCircle } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "../api/axios";

const EMAIL_REGEX = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/;
const PASSWORD_REGEX = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.{8,})/;
const LOGIN_URL = "/login";

const SignInForm = () => {
    const emailRef = useRef(null);
    const navigate = useNavigate();
    const location = useLocation();
    const from = location.state?.from?.pathname || "/";

    const [email, setEmail] = useState("");
    const [validEmail, setValidEmail] = useState(false);
    const [emailFocus, setEmailFocus] = useState(false);

    const [password, setPassword] = useState("");
    const [validPassword, setValidPassword] = useState(false);
    const [passwordFocus, setPasswordFocus] = useState(false);

    const [errMsg, setErrMsg] = useState("");

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
    }, [password]);

    useEffect(() => {
        setErrMsg("");
    }, [email, password]);

    const doLogin = async () => {
        const response = await axios.post(
            LOGIN_URL,
            JSON.stringify({email, password}),
            {
                headers: { "Content-Type": "application/json" },
                withCredentials: true
            }
        );

        return response?.data?.access_token;
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (EMAIL_REGEX.test(email) && PASSWORD_REGEX.test(password)) {
            try {
                await doLogin();

                setEmail("");
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
    };

    return (
        <section>
            <h1>Sign In</h1>
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
                <button
                    disabled={!validEmail || !validPassword}
                    type="submit"
                >
                    Sign In
                </button>
            </form>
            <p>
                Don&apos;t have an account yet?
                <Link to="/signup">Sign Up</Link>
            </p>
        </section>
    );
};

export default SignInForm;
