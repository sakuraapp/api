CREATE TABLE discriminators (
    id integer NOT NULL,
    name character varying(36) NOT NULL,
    value character varying(4) NOT NULL,
    owner_id integer
);

ALTER TABLE discriminators ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME discriminators_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

CREATE TABLE rooms (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    owner_id integer NOT NULL,
    private boolean NOT NULL
);

ALTER TABLE rooms ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME rooms_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

CREATE TABLE users (
    id integer NOT NULL,
    username character varying(36) NOT NULL,
    avatar text,
    provider character varying(255) NOT NULL,
    access_token text,
    refresh_token text,
    external_user_id character varying(255) NOT NULL,
    last_modified timestamp default now()
);

ALTER TABLE users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

CREATE TABLE user_roles(
     "id" INTEGER NOT NULL,
     "user_id" INTEGER NOT NULL,
     "room_id" INTEGER NOT NULL,
     "role_id" INTEGER NOT NULL
);

ALTER TABLE user_roles ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME user_roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


ALTER TABLE ONLY discriminators
    ADD CONSTRAINT discriminators_pkey PRIMARY KEY (id);

ALTER TABLE ONLY discriminators
    ADD CONSTRAINT discriminators_unique UNIQUE (name, value);


ALTER TABLE ONLY rooms
    ADD CONSTRAINT rooms_pkey PRIMARY KEY (id);


ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


ALTER TABLE ONLY user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);


CREATE INDEX discriminators_name_index ON discriminators USING btree (name);

CREATE INDEX discriminators_owner_id_index ON discriminators USING btree (owner_id);

CREATE INDEX rooms_owner_id_index ON rooms USING btree (owner_id);

CREATE INDEX user_roles_user_id_index ON user_roles USING btree (user_id);

CREATE INDEX user_roles_room_id_index ON user_roles USING btree (room_id);

ALTER TABLE ONLY discriminators
    ADD CONSTRAINT discriminators_owner_id_foreign FOREIGN KEY (owner_id) REFERENCES users(id);


ALTER TABLE ONLY rooms
    ADD CONSTRAINT rooms_owner_id_foreign FOREIGN KEY (owner_id) REFERENCES users(id);


ALTER TABLE user_roles
    ADD CONSTRAINT user_roles_user_id_foreign FOREIGN KEY(user_id) REFERENCES users(id);

ALTER TABLE user_roles
    ADD CONSTRAINT user_roles_room_id_foreign FOREIGN KEY(room_id) REFERENCES rooms(id);