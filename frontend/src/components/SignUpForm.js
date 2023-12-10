
import React, { useState } from "react";

const SignUpForm = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const handleEmailChange = (e) => {
        setEmail(e.target.value);
    };

    const handlePasswordChange = (e) => {
        setPassword(e.target.value);
    };

    const handleSubmit = (e) => {
        e.preventDefault();
    // Add your sign-in logic here
    };

    return (
        <div>
            <h2>Sign Up</h2>
            <form onSubmit={handleSubmit}>
                <div>
                    <label>Email:</label>
                    <input type="email" placeholder="Email" value={email} onChange={handleEmailChange} />
                </div>
                <div>
                    <label>Password:</label>
                    <input type="password" placeholder="Password" value={password} onChange={handlePasswordChange} />
                </div>
                <button type="submit">Sign In</button>
            </form>
        </div>
    );
};

export default SignUpForm;
