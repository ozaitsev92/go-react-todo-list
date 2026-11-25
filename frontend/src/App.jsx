import { Route, Routes } from "react-router-dom";
import { GuestOnly, RequireAuth } from "./auth/routes";
import Layout from "./components/Layout";
import NotFoundPage from "./pages/NotFoundPage";
import SignIn from "./pages/SignIn";
import SignUp from "./pages/SignUp";
import TodoList from "./pages/TodoList";

export default function App() {
    return (
        <Routes>
            <Route element={<Layout />}>
                <Route element={<GuestOnly />}>
                    <Route path="/signin" element={<SignIn />} />
                    <Route path="/signup" element={<SignUp />} />
                </Route>
                <Route element={<RequireAuth />}>
                    <Route index element={<TodoList />} />
                </Route>
                <Route path="*" element={<NotFoundPage />} />
            </Route>
        </Routes>
    );
}
