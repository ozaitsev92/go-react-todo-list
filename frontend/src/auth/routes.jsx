import { Navigate, Outlet, useLocation } from "react-router-dom";
import Alert from "react-bootstrap/Alert";
import Button from "react-bootstrap/Button";
import Spinner from "react-bootstrap/Spinner";
import useAuth from "./useAuth";

function Loading() {
    return (
        <div className="text-center my-5">
            <Spinner animation="border" role="status">
                <span className="visually-hidden">Loading…</span>
            </Spinner>
        </div>
    );
}

// RequireAuth renders child routes only when signed in; others go to /signin.
export function RequireAuth() {
    const { status, error, retry } = useAuth();
    const location = useLocation();

    if (status === "loading") {
        return <Loading />;
    }
    if (status === "error") {
        return (
            <Alert variant="danger" role="alert">
                {error}{" "}
                <Button variant="outline-danger" size="sm" onClick={retry}>
                    Retry
                </Button>
            </Alert>
        );
    }
    if (status !== "authenticated") {
        return <Navigate to="/signin" replace state={{ from: location }} />;
    }
    return <Outlet />;
}

// GuestOnly keeps signed-in users away from the sign-in and sign-up pages.
export function GuestOnly() {
    const { status } = useAuth();
    return status === "authenticated" ? <Navigate to="/" replace /> : <Outlet />;
}
