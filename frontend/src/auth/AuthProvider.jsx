import { createContext, useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import PropTypes from "prop-types";
import { authApi, errorMessage, isCancel, setUnauthorizedHandler } from "../lib/api";

export const AuthContext = createContext(null);

// status: "loading", then "authenticated", "anonymous" or "error" (check failed, not 401).
const LOADING = { status: "loading", user: null, error: null };
const ANONYMOUS = { status: "anonymous", user: null, error: null };

export function AuthProvider({ children }) {
    const [state, setState] = useState(LOADING);
    const navigate = useNavigate();
    const location = useLocation();
    const locationRef = useRef(location);
    useEffect(() => {
        locationRef.current = location;
    }, [location]);

    // Sets state only in callbacks (react-hooks/set-state-in-effect).
    const checkSession = useCallback(
        (signal) =>
            authApi.current(signal).then(
                (user) => setState({ status: "authenticated", user, error: null }),
                (error) => {
                    if (!isCancel(error)) {
                        setState(
                            error.response?.status === 401
                                ? ANONYMOUS
                                : { status: "error", user: null, error: errorMessage(error) },
                        );
                    }
                },
            ),
        [],
    );

    useEffect(() => {
        const controller = new AbortController();
        checkSession(controller.signal);
        return () => controller.abort();
    }, [checkSession]);

    // Any 401 from an authenticated call ends the session everywhere.
    useEffect(
        () =>
            setUnauthorizedHandler(() => {
                setState(ANONYMOUS);
                navigate("/signin", { replace: true, state: { from: locationRef.current } });
            }),
        [navigate],
    );

    const login = useCallback(async (email, password) => {
        const user = await authApi.login(email, password);
        setState({ status: "authenticated", user, error: null });
        return user;
    }, []);

    // Clears local state even if the request fails; the cookie then just expires.
    const logout = useCallback(async () => {
        try {
            await authApi.logout();
        } catch {
            // best effort
        } finally {
            setState(ANONYMOUS);
            navigate("/signin", { replace: true });
        }
    }, [navigate]);

    const value = useMemo(
        () => ({
            ...state,
            login,
            logout,
            retry: () => {
                setState(LOADING);
                checkSession();
            },
        }),
        [state, login, logout, checkSession],
    );

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

AuthProvider.propTypes = {
    children: PropTypes.node.isRequired,
};
