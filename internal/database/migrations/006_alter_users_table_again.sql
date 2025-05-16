ALTER TABLE users
ADD COLUMN is_admin BOOLEAN AFTER password;

-- Add default Admin Account for Debugging
INSERT INTO users (first_name, last_name, username, email, password, is_admin)
 VALUES ('Admin', '', 'Admin', 'admin@admin.com', '$2a$10$lYFnIxx/T20AdDYF5VFUAOrYKfQEwJxdWVocE4ccvs7C6wjvlgIle', 1);