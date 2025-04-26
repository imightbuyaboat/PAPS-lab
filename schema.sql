CREATE TABLE users (
    login TEXT PRIMARY KEY,
    hash TEXT NOT NULL,
    priveleged BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE register (
    id SERIAL PRIMARY KEY,
    organization VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
	CONSTRAINT phone_format CHECK (
        phone ~ '^\+7-\d{3}-\d{3}-\d{2}-\d{2}$'
    )
);