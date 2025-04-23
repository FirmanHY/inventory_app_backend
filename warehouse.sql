-- Setup Struktur Database Warehosue
CREATE DATABASE IF NOT EXISTS warehouse;
USE warehouse;

-- Tabel `users`
CREATE TABLE users (
    user_id char(36) PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    role ENUM('admin', 'warehouse_admin', 'warehouse_manager') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tabel `jenis_barang`
CREATE TABLE item_types (
    type_id char(36) PRIMARY KEY,
    type_name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tabel `satuan`
CREATE TABLE units (
    unit_id char(36) PRIMARY KEY,
    unit_name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tabel `barang`
CREATE TABLE items (
    item_id char(36) PRIMARY KEY,
    type_id char(36) NOT NULL,
    unit_id char(36) NOT NULL,
    item_name VARCHAR(255) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    image VARCHAR(255),
    minimum_stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (type_id) REFERENCES item_types(type_id) ON DELETE RESTRICT,
    FOREIGN KEY (unit_id) REFERENCES units(unit_id) ON DELETE RESTRICT
);

-- Tabel `transaksi`
CREATE TABLE transactions (
    transaction_id char(36) PRIMARY KEY,
    item_id char(36) NOT NULL,
    date DATE NOT NULL,
    quantity INT NOT NULL,
    transaction_type ENUM('in', 'out') NOT NULL,
    description TEXT,
    user_id char(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (item_id) REFERENCES items(item_id) ON DELETE RESTRICT,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE RESTRICT
);

DELIMITER $$

-- Trigger untuk transaksi masuk
CREATE TRIGGER stock_in
AFTER INSERT ON transactions
FOR EACH ROW
BEGIN
    IF NEW.transaction_type = 'in' THEN
        UPDATE items
        SET stock = stock + NEW.quantity
        WHERE item_id = NEW.item_id;
    END IF;
END $$

-- Trigger untuk transaksi keluar
CREATE TRIGGER stock_out
AFTER INSERT ON transactions
FOR EACH ROW
BEGIN
    IF NEW.transaction_type = 'out' THEN
        UPDATE items
        SET stock = stock - NEW.quantity
        WHERE item_id = NEW.item_id;
    END IF;
END $$

-- Trigger untuk menghapus transaksi
CREATE TRIGGER delete_stock
BEFORE DELETE ON transactions
FOR EACH ROW
BEGIN
    IF OLD.transaction_type = 'in' THEN
        UPDATE items
        SET stock = GREATEST(stock - OLD.quantity, 0)
        WHERE item_id = OLD.item_id;
    ELSEIF OLD.transaction_type = 'out' THEN
        UPDATE items
        SET stock = stock + OLD.quantity
        WHERE item_id = OLD.item_id;
    END IF;
END $$

DELIMITER ;

INSERT INTO users (user_id,username, password, full_name, role)
VALUES
    -- password: Admin1234
    (uuid(),'admin', '$2y$10$R3IqiZSysEAveFSBGBKlvuxfCZ3397ZkCWr.6aHgrboSST60zukpG', 'Administrator', 'admin'),
    -- password: Admingudang1234
    (uuid(),'warehouse_admin', '$2y$10$yQ1D9GSQ0dVOLeKdWgIQkuDbI2wKPZMTOSJUVibkHn7wopg3NTHHC', 'Warehouse Admin', 'warehouse_admin'),
    -- password: Kepalagudang1234
    (uuid(),'warehouse_manager', '$2y$10$ce1DNMyxFZQOosX/ZhEgJO9xoEp4V.PIIx25MGortv4S8mnghtW66', 'Warehouse Manager', 'warehouse_manager');

INSERT INTO item_types (type_id,type_name)
VALUES (uuid(),'Electronics'), (uuid(),'Clothing'), (uuid(),'Food');

INSERT INTO units (unit_id,unit_name)
VALUES (uuid(),'EA'), (uuid(),'KG'), (uuid(),'M');
