const mockTodoResponse = {
    data: [
        {
            id: 1,
            text: "Test task 1",
            completed: false,
            user_id: 1,
            isEditing: false
        },
        {
            id: 2,
            text: "Test task 2",
            completed: true,
            user_id: 1,
            isEditing: false
        }
    ]
};

const mockAuthResponse = {
    data: {
        accessToken: "test-access-token"
    }
};

export default {
    create(options = {}) {
        return {
            get: jest.fn().mockResolvedValue(mockTodoResponse),
            post: jest.fn().mockResolvedValue(mockAuthResponse),
            defaults: { headers: { common: {} } },
            ...options
        };
    }
};