import React from "react";
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import { Link } from "react-router-dom";

const NotFoundPage = () => {
    return (
        <Container className="text-center mt-5">
            <Row>
                <Col>
                    <h1>404 - Page Not Found</h1>
                </Col>
            </Row>
            <Row>
                <Col>
                    <p>The page you are looking for does not exist.</p>
                </Col>
            </Row>
            <Row>
                <Col>
                    <Link to="/">Go back to the homepage</Link>
                </Col>
            </Row>
        </Container>
    );
};

export default NotFoundPage;