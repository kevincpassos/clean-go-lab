CREATE TABLE companies (
    id              BIGSERIAL PRIMARY KEY,
    slug            VARCHAR(80)  NOT NULL,
    name            VARCHAR(120) NOT NULL,
    trade_name      VARCHAR(120),
    email           VARCHAR(255) NOT NULL,
    phone           VARCHAR(30),
    status          VARCHAR(20)  NOT NULL DEFAULT 'active',
    website         VARCHAR(255),
    photo_url       TEXT,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_companies_slug UNIQUE (slug),

    CONSTRAINT chk_companies_slug_length
        CHECK (char_length(BTRIM(slug)) BETWEEN 2 AND 80),

    CONSTRAINT chk_companies_name_length
        CHECK (char_length(BTRIM(name)) BETWEEN 2 AND 120),

    CONSTRAINT chk_companies_trade_name_length
        CHECK (
            trade_name IS NULL OR
            char_length(BTRIM(trade_name)) BETWEEN 2 AND 120
        ),

    CONSTRAINT chk_companies_email_length
        CHECK (char_length(BTRIM(email)) BETWEEN 5 AND 255),

    CONSTRAINT chk_companies_phone_length
        CHECK (
            phone IS NULL OR
            char_length(BTRIM(phone)) BETWEEN 8 AND 30
        ),

    CONSTRAINT chk_companies_website_length
        CHECK (
            website IS NULL OR
            char_length(BTRIM(website)) BETWEEN 8 AND 255
        ),

    CONSTRAINT chk_companies_photo_url_length
        CHECK (
            photo_url IS NULL OR
            char_length(BTRIM(photo_url)) BETWEEN 8 AND 2048
        ),

    CONSTRAINT chk_companies_status
        CHECK (status IN ('active', 'inactive', 'suspended'))
);

CREATE TABLE users (
    id                  BIGSERIAL PRIMARY KEY,
    name                VARCHAR(120) NOT NULL,
    email               VARCHAR(255) NOT NULL,
    password_hash       TEXT         NOT NULL,
    avatar_url          TEXT,
    phone               VARCHAR(30),
    email_verified_at   TIMESTAMPTZ,
    last_login_at       TIMESTAMPTZ,
    status              VARCHAR(20)  NOT NULL DEFAULT 'active',
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_users_email UNIQUE (email),

    CONSTRAINT chk_users_name_length
        CHECK (char_length(BTRIM(name)) BETWEEN 2 AND 120),

    CONSTRAINT chk_users_email_length
        CHECK (char_length(BTRIM(email)) BETWEEN 5 AND 255),

    CONSTRAINT chk_users_password_hash_length
        CHECK (char_length(BTRIM(password_hash)) >= 20),

    CONSTRAINT chk_users_avatar_url_length
        CHECK (
            avatar_url IS NULL OR
            char_length(BTRIM(avatar_url)) BETWEEN 8 AND 2048
        ),

    CONSTRAINT chk_users_phone_length
        CHECK (
            phone IS NULL OR
            char_length(BTRIM(phone)) BETWEEN 8 AND 30
        ),

    CONSTRAINT chk_users_status
        CHECK (status IN ('active', 'inactive', 'blocked')),

    CONSTRAINT chk_users_email_verified_after_created
        CHECK (
            email_verified_at IS NULL OR
            email_verified_at >= created_at
        ),

    CONSTRAINT chk_users_last_login_after_created
        CHECK (
            last_login_at IS NULL OR
            last_login_at >= created_at
        )
);

CREATE TABLE company_users (
    id              BIGSERIAL PRIMARY KEY,
    company_id      BIGINT       NOT NULL,
    user_id         BIGINT       NOT NULL,
    role            VARCHAR(30)  NOT NULL,
    is_owner        BOOLEAN      NOT NULL DEFAULT FALSE,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_company_users_company
        FOREIGN KEY (company_id)
        REFERENCES companies(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_company_users_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_company_users_company_user
        UNIQUE (company_id, user_id),

    CONSTRAINT chk_company_users_role_length
        CHECK (char_length(BTRIM(role)) BETWEEN 2 AND 30),

    CONSTRAINT chk_company_users_role
        CHECK (role IN ('owner', 'admin', 'manager', 'member', 'viewer'))
);



