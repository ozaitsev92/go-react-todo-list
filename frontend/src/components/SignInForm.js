
import React, { useRef, useState, useEffect } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { faCheck, faTimes } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "../api/axios";
import useInput from "../hooks/useInput";

const EMAIL_REGEX = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/;
const PASSWORD_REGEX = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.{8,})/;
const LOGIN_URL = "/login";

const SignInForm = () => {
    const emailRef = useRef(null);
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
                        required
                        {...emailAttrs}
                    />
                </div>
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
                    />
                </div>
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
