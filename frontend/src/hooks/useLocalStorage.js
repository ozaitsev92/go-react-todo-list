import { useState, useEffect } from "react";

const getLocalValue = (key, initValue) => {
    if (typeof window === "undefined") {
        return initValue;
    }

    const localValue = JSON.parse(window.localStorage.getItem(key));

    if (localValue) {
        return localValue;
    }

    if (initValue instanceof Function) {
        return initValue();
    }

    return initValue;
};

const useLocalStorage = (key, initValue) => {
    const [value, setValue] = useState(() => getLocalValue(key, initValue));

    useEffect(() => {
        if (typeof window !== "undefined") {
            window.localStorage.setItem("user", JSON.stringify(value));
        }
    }, [key, value]);

    return [value, setValue];
};

export default useLocalStorage;
