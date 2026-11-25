import Col from "react-bootstrap/Col";
import Row from "react-bootstrap/Row";
import SignInForm from "../components/SignInForm";

export default function SignIn() {
    return (
        <Row className="justify-content-center">
            <Col md={8} lg={6}>
                <SignInForm />
            </Col>
        </Row>
    );
}
