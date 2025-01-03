-- -- postgres syntax
-- CREATE TABLE users (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(150) NOT NULL, -- encrypted
--     email VARCHAR(100) NOT NULL UNIQUE,
--     password VARCHAR(255) NOT NULL,
--     role VARCHAR(50) NOT NULL,
--     id_card_number VARCHAR(50) NOT NULL UNIQUE, -- encrypted
--     id_family_card_number VARCHAR(50) NOT NULL, -- encrypted
--     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
-- );

INSERT INTO users (name, email, password, role, id_card_number, id_family_card_number)
VALUES ('John Doe', 'john.doe', 'password', 'admin', pgp_sym_encrypt( '7271011208920003', 'secret'), pgp_sym_encrypt( '7271011208920001', 'secret') );