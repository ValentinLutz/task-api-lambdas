CREATE TABLE IF NOT EXISTS public.tasks
(
    task_id         UUID    NOT NULL UNIQUE,
    title           TEXT,
    description     TEXT,
    PRIMARY KEY (task_id)
);