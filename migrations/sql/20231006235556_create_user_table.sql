-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.users
(
    user_id uuid NOT NULL,
    username varchar(15) NOT NULL,
    "password" text NOT NULL,
    email text NOT NULL,
    last_login_at timestamp NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    deleted_at timestamp NULL,
    is_email_verified bool NULL DEFAULT false,
    CONSTRAINT users_pkey PRIMARY KEY (user_id),
    CONSTRAINT unique_username UNIQUE (username), -- Tambahkan indeks unik pada kolom "username"
    CONSTRAINT unique_email UNIQUE (email) -- Tambahkan indeks unik pada kolom "email"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.users;
-- +goose StatementEnd
