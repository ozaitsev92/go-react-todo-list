import { useCallback, useEffect, useState } from "react";
import { errorMessage, isCancel, tasksApi } from "../lib/api";

// useTasks owns the task list. Mutations apply the server's response and
// resolve to true/false (the message is in `error`) instead of throwing.
export default function useTasks() {
    const [tasks, setTasks] = useState([]);
    const [status, setStatus] = useState("loading");
    const [error, setError] = useState(null);
    const [busyIds, setBusyIds] = useState(() => new Set());

    // Sets state only in callbacks (react-hooks/set-state-in-effect).
    const load = useCallback(
        (signal) =>
            tasksApi.list(signal).then(
                (list) => {
                    setTasks(list);
                    setError(null);
                    setStatus("ready");
                },
                (err) => {
                    if (!isCancel(err)) {
                        setError(errorMessage(err));
                        setStatus("error");
                    }
                },
            ),
        [],
    );

    useEffect(() => {
        const controller = new AbortController();
        load(controller.signal);
        return () => controller.abort();
    }, [load]);

    const reload = useCallback(() => {
        setStatus("loading");
        setError(null);
        load();
    }, [load]);

    const mutate = useCallback(async (id, request, apply) => {
        setError(null);
        if (id) {
            setBusyIds((prev) => new Set(prev).add(id));
        }
        try {
            const result = await request();
            setTasks((prev) => apply(prev, result));
            return true;
        } catch (err) {
            setError(errorMessage(err));
            return false;
        } finally {
            if (id) {
                setBusyIds((prev) => {
                    const next = new Set(prev);
                    next.delete(id);
                    return next;
                });
            }
        }
    }, []);

    const replace = (prev, updated) => prev.map((t) => (t.id === updated.id ? updated : t));

    const create = useCallback(
        (text) =>
            mutate(
                null,
                () => tasksApi.create(text),
                (prev, created) => [...prev, created],
            ),
        [mutate],
    );

    const update = useCallback(
        (id, text) => mutate(id, () => tasksApi.update(id, text), replace),
        [mutate],
    );

    const setCompleted = useCallback(
        (id, completed) => mutate(id, () => tasksApi.setCompleted(id, completed), replace),
        [mutate],
    );

    const remove = useCallback(
        (id) =>
            mutate(
                id,
                () => tasksApi.remove(id),
                (prev) => prev.filter((t) => t.id !== id),
            ),
        [mutate],
    );

    return {
        tasks,
        status,
        error,
        busyIds,
        reload,
        clearError: () => setError(null),
        create,
        update,
        setCompleted,
        remove,
    };
}
