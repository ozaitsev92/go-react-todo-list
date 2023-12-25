import React from "react";
import Container from "react-bootstrap/Container";
import SignInForm from "../components/SignInForm";

const SignIn = () => {
    return (
        <Container className="mt-5">
            <SignInForm />
        </Container>
    );
};

export default SignIn;