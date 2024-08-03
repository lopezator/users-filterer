CREATE DATABASE users;
ALTER DATABASE users OWNER TO root;
CREATE TABLE public.users (
                              id
                                  INT8 DEFAULT unique_rowid() NOT NULL PRIMARY KEY,
                              display_name
                                  VARCHAR(50) NOT NULL,
                              email
                                  VARCHAR(100) NOT NULL UNIQUE,
                              created_at
                                  TIMESTAMP DEFAULT now()
);
ALTER TABLE public.users OWNER TO root;
INSERT
INTO
    public.users (id, display_name, email, created_at)
VALUES
    (
        991650920442527745,
        'alice',
        'alice@example.com',
        '2024-08-03 15:17:24.365002'
    );
INSERT
INTO
    public.users (id, display_name, email, created_at)
VALUES
    (
        991650920442626049,
        'bob',
        'bob@example.com',
        '2024-08-03 15:17:24.365002'
    );
INSERT
INTO
    public.users (id, display_name, email, created_at)
VALUES
    (
        991650920442658817,
        'charlie',
        'charlie@example.com',
        '2024-08-03 15:17:24.365002'
    );
