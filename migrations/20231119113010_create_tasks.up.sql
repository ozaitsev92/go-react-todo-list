CREATE TABLE tasks (
    id bigserial not null primary key,
    task_text varchar not null,
    task_order smallint not null,
    is_done boolean not null,
    user_id bigint not null,
    CONSTRAINT fk_user_id
      FOREIGN KEY(user_id)
      REFERENCES users(id)
);