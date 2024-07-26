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
    username VARCHAR(100) NOT NULL UNIQUE,
    nombre VARCHAR(100),
    apellido VARCHAR(100),
    doc VARCHAR(100),
    tipo_doc ENUM('DNI', 'PASAPORTE', 'CEDULA'),
    email VARCHAR(100) NOT NULL UNIQUE,
    fecha_nacimiento DATETIME,
    activo BOOLEAN
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

CREATE TABLE IF NOT EXISTS auditoria (
        id INT PRIMARY KEY AUTO_INCREMENT,
        id_usuario INT NOT NULL,
        id_cotizacion INT NOT NULL,
        accion VARCHAR(255),
        columna_afectada VARCHAR(255),
        viejo_valor VARCHAR(255),
        nuevo_valor VARCHAR(255),
        fecha DATETIME,
        FOREIGN KEY (id_usuario) REFERENCES usuario(id),
        FOREIGN KEY (id_cotizacion) REFERENCES cotizacion(id)
    );