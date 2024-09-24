create extension if not exists "uuid-ossp";
-- public.user_sessions definition

-- Drop table

-- DROP TABLE public.user_sessions;

CREATE TABLE public.user_sessions (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	auth_token text NULL,
	refresh_token_id text NULL,
	notification_token text NULL,
	expired_at timestamptz NULL,
	status text NULL,
	user_id text NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	created_by text NULL,
	updated_by text NULL,
	record_flag text NULL,
	CONSTRAINT uni_user_sessions_id PRIMARY KEY (id)
);


-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	"role" text NULL,
	"name" text NULL,
	email text NULL,
	"password" text NULL,
	email_verified_at timestamptz NULL,
	last_login_method text NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	created_by text NULL,
	updated_by text NULL,
	record_flag text NULL,
	CONSTRAINT uni_users_id PRIMARY KEY (id)
);

INSERT INTO public.users

(id, "role", "name", email, "password", email_verified_at, last_login_method, created_at, updated_at, created_by, updated_by, record_flag)
VALUES('8ccb3f39-a6c3-4f54-aee2-3c641a0c5fc0'::uuid, 'superadmin', 'fajar-edist', 'aegis@aegis.com', '$2a$10$9c5ua1JVAikmlyNuD3744eqKWem22ANMaERuMmh0uukZJQL3ZE.Bq', NULL, '', '2024-09-24 22:37:45.265', '2024-09-24 22:38:39.486', 'service.(*userService).Register', '', 'ACTIVE');
INSERT INTO public.users
(id, "role", "name", email, "password", email_verified_at, last_login_method, created_at, updated_at, created_by, updated_by, record_flag)
VALUES('ed12d327-5bd7-42c1-8908-6d168b688660'::uuid, 'user', 'fajar-edist', 'user3@aegis.com', '$2a$10$9c5ua1JVAikmlyNuD3744eqKWem22ANMaERuMmh0uukZJQL3ZE.Bq', NULL, '', '2024-09-24 22:37:45.265', '2024-09-24 22:38:39.486', 'service.(*userService).Register', '', 'ACTIVE');
INSERT INTO public.users
(id, "role", "name", email, "password", email_verified_at, last_login_method, created_at, updated_at, created_by, updated_by, record_flag)
VALUES('5b5f3cda-bd77-4470-a9cc-7056783fa2a3'::uuid, 'superadmin', 'fajar-edist', 'fajar@aegis.com', '$2a$10$9c5ua1JVAikmlyNuD3744eqKWem22ANMaERuMmh0uukZJQL3ZE.Bq', NULL, '', '2024-09-24 22:37:45.265', '2024-09-24 22:44:34.012', 'service.(*userService).Register', '', 'ACTIVE');
INSERT INTO public.users
(id, "role", "name", email, "password", email_verified_at, last_login_method, created_at, updated_at, created_by, updated_by, record_flag)
VALUES('6892063d-38bf-4ee6-a774-29b2b9a6aa77'::uuid, 'user', 'fajar-edist', 'user5@aegis.com', '$2a$10$9c5ua1JVAikmlyNuD3744eqKWem22ANMaERuMmh0uukZJQL3ZE.Bq', NULL, '', '2024-09-24 22:37:45.265', '2024-09-24 22:38:39.486', 'service.(*userService).Register', '', 'ACTIVE');
INSERT INTO public.users
(id, "role", "name", email, "password", email_verified_at, last_login_method, created_at, updated_at, created_by, updated_by, record_flag)
VALUES('e8e1769a-d8ec-49a1-8f3d-5cbc6bfaf71c'::uuid, 'user', 'fajar-edista', 'user4@aegis.com', '$2a$10$9c5ua1JVAikmlyNuD3744eqKWem22ANMaERuMmh0uukZJQL3ZE.Bq', NULL, '', '2024-09-24 22:37:45.265', '2024-09-24 22:44:53.526', 'service.(*userService).Register', '', 'ACTIVE');
INSERT INTO public.users
(id, "role", "name", email, "password", email_verified_at, last_login_method, created_at, updated_at, created_by, updated_by, record_flag)
VALUES('7a70e246-a9f7-44d4-ac4f-de92feef3970'::uuid, 'user', 'fajar-edista', 'user2@aegis.com', '$2a$10$9c5ua1JVAikmlyNuD3744eqKWem22ANMaERuMmh0uukZJQL3ZE.Bq', NULL, '', '2024-09-24 22:37:45.265', '2024-09-24 22:45:50.517', 'service.(*userService).Register', '', 'ACTIVE');
