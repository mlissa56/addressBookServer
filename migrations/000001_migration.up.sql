CREATE TABLE IF NOT EXISTS address_book (
    id SERIAL PRIMARY KEY,
    name TEXT,
    last_name TEXT,
    middle_name TEXT,
    address TEXT,
    phone TEXT
); 
