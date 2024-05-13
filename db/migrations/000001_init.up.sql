CREATE TABLE IF NOT EXISTS accounts (
	id serial NOT NULL,
	owner varchar(255) NOT NULL,
	balance bigint NOT NULL,
	currency varchar(255) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	PRIMARY KEY(id)
);
CREATE INDEX owner_index ON accounts (owner);

CREATE TABLE IF NOT EXISTS entries (
	id serial NOT NULL,
	account_id int NOT NULL,
	amount bigint NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	PRIMARY KEY(id)
);
CREATE INDEX account_id_index ON entries (account_id);

CREATE TABLE IF NOT EXISTS transfers (
	id serial NOT NULL,
	from_account_id int NOT NULL,
	to_account_id int NOT NULL,
	amount bigint NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	PRIMARY KEY(id)
);

CREATE INDEX from_account_id_index ON transfers (from_account_id);
CREATE INDEX to_account_id_index ON transfers (to_account_id);
CREATE INDEX from_account_id_to_account_id_index ON transfers (from_account_id, to_account_id);

ALTER TABLE entries
ADD FOREIGN KEY (account_id) REFERENCES accounts(id) ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE transfers
ADD FOREIGN KEY(from_account_id) REFERENCES accounts(id) ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE transfers
ADD FOREIGN KEY(to_account_id) REFERENCES accounts(id) ON UPDATE NO ACTION ON DELETE NO ACTION;