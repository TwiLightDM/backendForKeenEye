WITH new_user AS (
    INSERT INTO users (login, password, salt, role)
        VALUES ('admin', '$2a$04$llb7M3x.Y3GV8axQnxc/..X7NHunBwT5fVx1nQtkSMzKIgNv86p1W', 'X2fzSkHued', 'admin')
        RETURNING id
)
INSERT INTO admins (id)
SELECT id FROM new_user;