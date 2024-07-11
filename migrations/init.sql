-- Crear el esquema 'go'
CREATE DATABASE IF NOT EXISTS `go`;

-- Usar el esquema 'go'
USE `go`;

-- Crear la tabla 'criptomoneda'
CREATE TABLE IF NOT EXISTS criptomoneda (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL
);

-- Insertar datos iniciales en 'criptomoneda'
INSERT INTO criptomoneda (nombre) VALUES
('Bitcoin'),
('Ethereum'),
('Litecoin');

-- Crear la tabla 'cotizacion'
CREATE TABLE IF NOT EXISTS cotizacion (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_criptomoneda INT,
    fecha DATETIME,
    valor DECIMAL(18, 2),
    FOREIGN KEY (id_criptomoneda) REFERENCES criptomoneda(id)
);

-- Insertar datos iniciales en 'cotizacion'
INSERT INTO cotizacion (id_criptomoneda, fecha, valor) VALUES
(1, '2024-07-01 08:00:00', 35000.00),
(2, '2024-07-01 12:30:00', 2300.50),
(3, '2024-07-01 15:45:00', 150.75),
(1, '2024-07-02 09:15:00', 35500.00),
(2, '2024-07-02 10:00:00', 2400.25),
(3, '2024-07-02 14:20:00', 155.80);
