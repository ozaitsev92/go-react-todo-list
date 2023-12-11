import React, { useEffect, useState } from "react";
import { useLocation, Navigate, Outlet } from "react-router-dom";
import useAuth from "../hooks/useAuth";
import axios from "../api/axios";

const CURRENT_USER_URL = "/users-current";

const RequireAuth = () => {
    const { auth, setAuth } = useAuth();
    const location = useLocation();
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const getCurrentUser = async () => {
            try {
                const response = await axios.get(
                    CURRENT_USER_URL,
                    {
                        headers: { "Content-Type": "application/json" },
                        withCredentials: true
                    }
                );

                const user = response?.data;
                setAuth({user});
                setIsLoading(false);
            } catch (error) {
                setIsLoading(false);
            }
        };

        if (!auth.user && location.pathname !== "/signin") {
            getCurrentUser();
        } else {
            setIsLoading(false);
        }

        // eslint-disable-next-line
    }, [auth.user]);

    return (
        <>
            {isLoading ? (
                <div>Loading...</div>
            ) : (
                !auth.user && location.pathname !== "/signin" ? (
                    <Navigate to="/signin" state={{from: location}} replace />
                ) : (
                    <Outlet />
                )
            )}
        </>
    );
};

export default RequireAuth;