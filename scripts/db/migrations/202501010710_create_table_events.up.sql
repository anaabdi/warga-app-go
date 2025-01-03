-- postgres syntax
-- requesting a changes which has to attach some documents to the request
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    event_type VARCHAR(50) NOT NULL, -- registrasi peristiwa pendudukan
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE event_requests (
    id SERIAL PRIMARY KEY,
    user_event_id INT NOT NULL,
    request_type VARCHAR(50) NOT NULL, -- kk, ktp, kia, 
    request_reason VARCHAR(255) NOT NULL, -- perubahan alamat, perubahan status, dll
    status VARCHAR(50) NOT NULL, -- pending, approved, rejected
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE event_docs (
    id SERIAL PRIMARY KEY,
    user_event_id INT NOT NULL,
    doc_type VARCHAR(50) NOT NULL, -- surat pengantar, kartu keluarga, akta kelahiran, akta nikah, dll
    doc_name VARCHAR(100) NOT NULL,
    doc_url VARCHAR(255) NOT NULL, -- https://s3.aws.cloud/lokasidirektori/filea.png -> /lokasidirektori/filea.png
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);