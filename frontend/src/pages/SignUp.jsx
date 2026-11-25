import Col from "react-bootstrap/Col";
import Row from "react-bootstrap/Row";
import SignUpForm from "../components/SignUpForm";

export default function SignUp() {
    return (
        <Row className="justify-content-center">
            <Col md={8} lg={6}>
                <SignUpForm />
            </Col>
        </Row>
    );
}
