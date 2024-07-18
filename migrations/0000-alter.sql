USE go;

ALTER TABLE criptomoneda
ADD COLUMN simbolo VARCHAR(5);

ALTER TABLA cotizacion
ADD COLUMN api VARCHAR(20);
