import { useEffect } from "react";
import { render } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, useLocation } from "react-router-dom";
import App from "../App";
import { AuthProvider } from "../auth/AuthProvider";

let currentPath = "";

function LocationProbe() {
    const { pathname } = useLocation();
    useEffect(() => {
        currentPath = pathname;
    }, [pathname]);
    return null;
}

// renderApp renders the real route tree at `route`, as the browser would.
export function renderApp(route = "/", state = undefined) {
    const user = userEvent.setup();
    const utils = render(
        <MemoryRouter initialEntries={[{ pathname: route, state }]}>
            <AuthProvider>
                <App />
                <LocationProbe />
            </AuthProvider>
        </MemoryRouter>,
    );
    return { user, ...utils, path: () => currentPath };
}
