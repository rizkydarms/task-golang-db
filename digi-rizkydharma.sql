--add foreign key
ALTER TABLE public.accounts ADD FOREIGN KEY (referral_account_id) REFERENCES public.auths(account_id);

ALTER TABLE public.transaction ADD FOREIGN KEY (transaction_category_id) 
REFERENCES public.transaction_categories(transaction_category_id);
-- 

INSERT INTO public.auths 
(username, "password", account_id)
VALUES('Rizky', 'Rizky123', 10);

-- Buat Query Insert 5 data accounts
INSERT INTO public.accounts
(name, balance, referral_account_id)
values
('Rizky', 500, 0),
('Dharma', 8750, 0),
('RizkyD', 7750, 0),
('RizkyDS', 9210, 0),
('RizkyDSS', 1550, 0);

-- Buat Query Insert 2 data transaction_categories
INSERT INTO public.transaction_categories
("name")
VALUES('Sembako'), ('Bukan Sembako');

-- Buat Query Insert 12 data transaction (kasih 1 transaksi tiap bulan)
INSERT INTO public."transaction"
(transaction_category_id, account_id, from_account_id, to_account_id, amount, transaction_date)
VALUES
(3, 123, 123, 234, 1000, now()),
(3, 123, 123, 234, 1100, now()),
(4, 123, 123, 234, 1200, now()),
(4, 123, 123, 234, 1300, now()),
(3, 123, 123, 234, 1400, now()),
(4, 123, 123, 234, 1500, now()),
(3, 123, 123, 234, 1600, now()),
(4, 123, 123, 234, 1700, now()),
(3, 123, 123, 234, 1800, now()),
(3, 123, 123, 234, 1900, now()),
(3, 123, 123, 234, 2000, now());

-- Buat Update Query
--Update nama di accounts (given account_id)
--Update balance di accounts (given account_id)
UPDATE public.accounts
SET "name"='RDSX', balance=500000, referral_account_id=0
WHERE account_id=12;

--Buat Query Select
--List semua data accounts
SELECT account_id, "name", balance, referral_account_id
FROM public.accounts;

--List semua data transactions join dengan accounts (account_id = account_id) dan tampilkan nama dari accounts
select * from accounts
join transaction on accounts.account_id = transaction.account_id 

--Query 1 data accounts dengan balance terbanyak
SELECT * FROM accounts ORDER BY balance desc LIMIT 1;

--Query semua transaction yg terjadi di bulan Mei (Bulan 5)
SELECT *
FROM transaction
WHERE transaction_date BETWEEN '2023-05-01' AND '2023-05-31';


--DDL

-- public.accounts definition

-- Drop table

-- DROP TABLE public.accounts;

CREATE TABLE public.accounts (
	account_id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
	"name" varchar NOT NULL,
	balance int8 DEFAULT 0 NOT NULL,
	referral_account_id int8 NULL,
	CONSTRAINT accounts_pk PRIMARY KEY (account_id)
);


-- public.accounts foreign keys

ALTER TABLE public.accounts ADD CONSTRAINT accounts_referral_account_id_fkey FOREIGN KEY (referral_account_id) REFERENCES public.auths(account_id);





-- public.auths definition

-- Drop table

-- DROP TABLE public.auths;

CREATE TABLE public.auths (
	auth_id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
	account_id int8 NOT NULL,
	username varchar(60) NOT NULL,
	"password" varchar NOT NULL,
	CONSTRAINT auths_pk PRIMARY KEY (auth_id),
	CONSTRAINT auths_unique UNIQUE (account_id),
	CONSTRAINT auths_unique_1 UNIQUE (username)
);





-- public.transaction_categories definition

-- Drop table

-- DROP TABLE public.transaction_categories;

CREATE TABLE public.transaction_categories (
	transaction_category_id int4 GENERATED ALWAYS AS IDENTITY NOT NULL,
	"name" varchar NULL,
	CONSTRAINT transaction_categories_pk PRIMARY KEY (transaction_category_id)
);




-- public."transaction" definition

-- Drop table

-- DROP TABLE public."transaction";

CREATE TABLE public."transaction" (
	transaction_id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
	transaction_category_id int8 NULL,
	account_id int8 NULL,
	from_account_id int8 NULL,
	to_account_id int8 NULL,
	amount int8 NULL,
	transaction_date timestamp NULL,
	CONSTRAINT transaction_pk PRIMARY KEY (transaction_id)
);


-- public."transaction" foreign keys

ALTER TABLE public."transaction" ADD CONSTRAINT transaction_transaction_category_id_fkey FOREIGN KEY (transaction_category_id) REFERENCES public.transaction_categories(transaction_category_id);

