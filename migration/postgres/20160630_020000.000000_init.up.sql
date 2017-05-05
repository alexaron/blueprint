SET TIME ZONE 'UTC';
-- Uncomment if a trigger below fails with an error.
-- CREATE LANGUAGE plpgsql;

-- ON UPDATE CURRENT_TIMESTAMP() MySQL analogue.
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
NEW.updated_at = now();
RETURN NEW;
END;
$$ language 'plpgsql';

-- ******************************************************************************
-- Create tables
-- ******************************************************************************
CREATE TABLE user_status (
    id BIGSERIAL PRIMARY KEY NOT NULL,

    status VARCHAR(25) NOT NULL,

    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);
CREATE TRIGGER update_user_status_updated_at BEFORE UPDATE
    ON user_status FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE "user" (
    id BIGSERIAL PRIMARY KEY NOT NULL,

    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password CHAR(60) NOT NULL,
    verification_code VARCHAR(128) NOT NULL UNIQUE,
    verified BOOL NOT NULL DEFAULT false,

    status_id BIGINT NOT NULL DEFAULT 1,

    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    CONSTRAINT f_user_status FOREIGN KEY (status_id) REFERENCES user_status (id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE TRIGGER update_user_updated_at BEFORE UPDATE
    ON "user" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

INSERT INTO user_status (id, status, created_at, updated_at, deleted_at) VALUES
    (1, 'active',   CURRENT_TIMESTAMP,  NULL,  NULL),
    (2, 'inactive', CURRENT_TIMESTAMP,  NULL,  NULL);

CREATE TABLE note (
    id BIGSERIAL PRIMARY KEY NOT NULL,

    name TEXT NOT NULL,

    user_id BIGINT NOT NULL,

    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    CONSTRAINT f_note_user FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE TRIGGER update_note_updated_at BEFORE UPDATE
    ON "note" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
