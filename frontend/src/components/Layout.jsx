import { Link, Outlet } from "react-router-dom";
import Button from "react-bootstrap/Button";
import Container from "react-bootstrap/Container";
import Navbar from "react-bootstrap/Navbar";
import useAuth from "../auth/useAuth";

export default function Layout() {
    const { user, logout } = useAuth();

    return (
        <>
            <Navbar bg="light" className="mb-4">
                <Container>
                    <Navbar.Brand as={Link} to="/">
                        TODO App
                    </Navbar.Brand>
                    {user && (
                        <div className="d-flex align-items-center gap-3">
                            <span className="text-muted">{user.email}</span>
                            <Button variant="outline-dark" size="sm" onClick={logout}>
                                Log out
                            </Button>
                        </div>
                    )}
                </Container>
            </Navbar>
            <Container as="main">
                <Outlet />
            </Container>
        </>
    );
}
