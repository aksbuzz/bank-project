CREATE TABLE IF NOT EXISTS books (
	id serial NOT NULL,
	title varchar(255) NOT NULL,
	author varchar(255) NOT NULL,
	year int NOT NULL,
	quantity int NOT NULL CHECK (quantity >= 0),
	PRIMARY KEY(id)
);
CREATE INDEX author_index ON books (author);

CREATE TYPE enum_membership_type AS ENUM ('basic', 'premium');

CREATE TABLE IF NOT EXISTS members (
	id serial NOT NULL,
	name varchar(255) NOT NULL,
	email varchar(255) NOT NULL UNIQUE,
	phone varchar(20) NOT NULL,
	membership_type enum_membership_type NOT NULL DEFAULT 'basic',
	PRIMARY KEY(id)
);
CREATE INDEX membership_type_index ON members (membership_type);

CREATE TABLE IF NOT EXISTS loans (
	id serial NOT NULL,
	book_id int NOT NULL,
	member_id int NOT NULL,
	loan_date DATE NOT NULL DEFAULT CURRENT_DATE,
	due_date DATE NOT NULL DEFAULT CURRENT_DATE + INTERVAL '30 days',
	return_date DATE,
	overdue_fee DECIMAL(5, 2),
	PRIMARY KEY(id)
);
CREATE INDEX book_id_index ON loans (book_id);
CREATE INDEX member_id_index ON loans (member_id);

ALTER TABLE loans
ADD FOREIGN KEY (book_id) REFERENCES books(id) ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE loans
ADD FOREIGN KEY(member_id) REFERENCES members(id) ON UPDATE NO ACTION ON DELETE NO ACTION;
