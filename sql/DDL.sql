CREATE TABLE mst_customer (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    address VARCHAR(255)
);

CREATE TABLE mst_employee (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    address VARCHAR(255)
);

CREATE TABLE mst_product (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    unit VARCHAR(7) NOT NULL,
    price INT NOT NULL
);

CREATE TABLE trx_bill (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    employee_id INT,
    bill_date DATE NOT NULL,
    entry_date DATE NOT NULL,
    finish_date DATE NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES mst_customer (id) ON DELETE CASCADE,
    FOREIGN KEY (employee_id) REFERENCES mst_employee (id) ON DELETE SET NULL
);

CREATE TABLE trx_bill_detail (
    id SERIAL PRIMARY KEY,
    bill_id INT NOT NULL,
    product_id INT NOT NULL,
    qty INT NOT NULL,
    product_price INT NOT NULL,
    FOREIGN KEY (bill_id) REFERENCES trx_bill (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES mst_product (id) ON DELETE CASCADE
);
