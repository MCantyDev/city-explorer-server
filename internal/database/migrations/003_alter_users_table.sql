ALTER TABLE users
ADD COLUMN email VARCHAR(255) NOT NULL UNIQUE AFTER username,
ADD COLUMN first_name varchar(100) NOT NULL AFTER id,
ADD COLUMN last_name varchar(100) NOT NULL AFTER first_name;