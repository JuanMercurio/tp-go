-- Crear el esquema 'go'
CREATE DATABASE IF NOT EXISTS `go`;

-- Usar el esquema 'go'
USE `go`;

-- Crear la tabla 'criptomoneda'
CREATE TABLE IF NOT EXISTS criptomoneda (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    simbolo VARCHAR(10) NOT NULL
);

-- Crear la tabla 'usuario'
CREATE TABLE IF NOT EXISTS usuario (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE,
    activo BOOLEAN,
);

-- Crear la tabla 'usuario_criptomoneda' para la relación entre 'usuario' y 'criptomoneda'
CREATE TABLE IF NOT EXISTS usuario_criptomoneda (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_usuario INT,
    id_criptomoneda INT,
    FOREIGN KEY (id_usuario) REFERENCES usuario(id),
    FOREIGN KEY (id_criptomoneda) REFERENCES criptomoneda(id)
);

-- Crear la tabla 'cotizacion'
CREATE TABLE IF NOT EXISTS cotizacion (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_criptomoneda INT,
    fecha DATETIME,
    valor DECIMAL(18, 2),
    api VARCHAR(50), 
    FOREIGN KEY (id_criptomoneda) REFERENCES criptomoneda(id)
);
