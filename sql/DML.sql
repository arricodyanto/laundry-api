INSERT INTO mst_customer (id, name, contact, address, is_employee) VALUES
(1, 'Mirna', '085664876443', 'Jakarta', true),
(2, 'Jessica', '0812654987', 'Bandung', false);

INSERT INTO mst_product (id, name, unit, price) VALUES
(1, 'Cuci + Setrika', 'KG', 7000),
(2, 'Laundry Bedcover', 'Buah', 50000),
(3, 'Laundry Boneka', 'Buah', 25000);

INSERT INTO trx_bill (id, customer_id, employee_id, bill_date, entry_date, finish_date) VALUES
(1, 2, 1, '2022-08-18', '2022-08-18', '2022-08-20');

INSERT INTO trx_bill_detail (id, bill_id, product_id, qty, product_price) VALUES
(1, 1, 1, 5, 35000),
(2, 1, 2, 1, 50000),
(3, 1, 3, 2, 50000);