CREATE TABLE region (
    id           CHAR(2)  PRIMARY KEY, -- Código ISO-3166-2 (AP, TA, …)
    number       SMALLINT NOT NULL,    -- 1–16
    roman_number TEXT     NOT NULL,
    name         TEXT     NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);