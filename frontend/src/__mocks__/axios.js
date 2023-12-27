const mockResponse = {
    data: [
        {
            id: 1,
            taskText: "Test task 1",
            isDone: false,
            userId: 1,
            isEditing: false,
            taskOrder: 1,
        },
        {
            id: 2,
            taskText: "Test task 2",
            isDone: true,
            userId: 1,
            isEditing: false,
            taskOrder: 2,
        },
    ]
}

export default {
    create(options = {}) {
        return {
            get: jest.fn().mockResolvedValue(mockResponse),
            defaults: { headers: { common: {} } },
            ...options,
        };
    },
};