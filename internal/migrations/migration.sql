CREATE DATABASE balance;
USE balance;
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
);

CREATE TABLE IF NOT EXISTS balance_service (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,
    balance DECIMAL(20, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO users (username, email) VALUES 
('Marsel', 'Marsel@mail.com'),
('Nastya', 'Nastya@mail.com'),
('Robert', 'Robert@mail.com');

INSERT INTO balance_service (user_id, balance) VALUES 
(1, 1000.00), 
(2, 1500.00), 
(3, 2000.00); 

