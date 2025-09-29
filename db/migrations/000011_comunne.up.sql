CREATE TABLE commune (
    id        VARCHAR(100) PRIMARY KEY, 
    region_id CHAR(2)  NOT NULL,
    name      TEXT     NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_region
        FOREIGN KEY (region_id) REFERENCES region(id)
);