--------------------------------------------------------------------------
-- Testdaten für die Tabelle "user"
--------------------------------------------------------------------------
INSERT INTO klcs."user" (username, name) VALUES
('johndoe', 'John Doe'),
('janedoe', 'Jane Doe'),
('alexsmith', 'Alex Smith');

--------------------------------------------------------------------------
-- Testdaten für die Tabelle "role"
--------------------------------------------------------------------------
INSERT INTO klcs.role (name) VALUES
('ADMIN'),
('UHD'),
('SELLER');

--------------------------------------------------------------------------
-- Testdaten für die Tabelle "shop"
--------------------------------------------------------------------------
INSERT INTO klcs."shop" (name) VALUES
('Shop A'),
('Shop B'),
('Shop C');

--------------------------------------------------------------------------
-- Testdaten für die Tabelle "account"
--------------------------------------------------------------------------
INSERT INTO klcs.account (locked, holder_name) VALUES
(false, 'John Doe'),
(true, 'Jane Doe'),
(false, 'Alex Smith');

--------------------------------------------------------------------------
-- Testdaten für die Tabelle "printer"
--------------------------------------------------------------------------
INSERT INTO klcs.printer (name, shop_id) VALUES
('Printer A', (SELECT id FROM klcs."shop" WHERE name = 'Shop A')),
('Printer B', (SELECT id FROM klcs."shop" WHERE name = 'Shop B')),
('Printer C', (SELECT id FROM klcs."shop" WHERE name = 'Shop C'));

--------------------------------------------------------------------------
-- Testdaten für die Tabelle "article"
--------------------------------------------------------------------------
INSERT INTO klcs.article (name, description, price, category, stock_amount, shop_id, printer_id) VALUES
('Article 1', 'Description of Article 1', 1500, 'Category A', 100, (SELECT id FROM klcs."shop" WHERE name = 'Shop A'), (SELECT id FROM klcs.printer WHERE name = 'Printer A')),
('Article 2', 'Description of Article 2', 2500, 'Category B', 200, (SELECT id FROM klcs."shop" WHERE name = 'Shop B'), (SELECT id FROM klcs.printer WHERE name = 'Printer B')),
('Article 3', 'Description of Article 3', 3500, 'Category C', 300, (SELECT id FROM klcs."shop" WHERE name = 'Shop C'), (SELECT id FROM klcs.printer WHERE name = 'Printer C'));

--------------------------------------------------------------------------
-- Testdaten für die Tabelle "transaction"
--------------------------------------------------------------------------
INSERT INTO klcs.transaction (type, amount, description, account_id) VALUES
('CASH', 1500, 'Cash transaction for Article 1', NULL),
('CARD', 2500, 'Card transaction for Article 2', (SELECT id FROM klcs.account WHERE holder_name = 'John Doe')),
('CARD', 3500, 'Card transaction for Article 3', (SELECT id FROM klcs.account WHERE holder_name = 'Jane Doe'));

--------------------------------------------------------------------------
-- Testdaten für die Tabelle "user_shop_role"
--------------------------------------------------------------------------
INSERT INTO klcs.user_shop_role (user_id, shop_id, role_id) VALUES
((SELECT id FROM klcs."user" WHERE username = 'johndoe'), (SELECT id FROM klcs."shop" WHERE name = 'Shop A'), (SELECT id FROM klcs.role WHERE name = 'ADMIN')),
((SELECT id FROM klcs."user" WHERE username = 'janedoe'), (SELECT id FROM klcs."shop" WHERE name = 'Shop B'), (SELECT id FROM klcs.role WHERE name = 'UHD')),
((SELECT id FROM klcs."user" WHERE username = 'alexsmith'), (SELECT id FROM klcs."shop" WHERE name = 'Shop C'), (SELECT id FROM klcs.role WHERE name = 'SELLER'));

--------------------------------------------------------------------------
-- Testdaten für die Tabelle "user_article_transaction"
--------------------------------------------------------------------------
INSERT INTO klcs.user_article_transaction (user_id, article_id, transaction_id, pieces) VALUES
((SELECT id FROM klcs."user" WHERE username = 'johndoe'), (SELECT id FROM klcs.article WHERE name = 'Article 1'), (SELECT id FROM klcs.transaction WHERE description = 'Cash transaction for Article 1'), 1),
((SELECT id FROM klcs."user" WHERE username = 'janedoe'), (SELECT id FROM klcs.article WHERE name = 'Article 2'), (SELECT id FROM klcs.transaction WHERE description = 'Card transaction for Article 2'), 2),
((SELECT id FROM klcs."user" WHERE username = 'alexsmith'), (SELECT id FROM klcs.article WHERE name = 'Article 3'), (SELECT id FROM klcs.transaction WHERE description = 'Card transaction for Article 3'), 3);
