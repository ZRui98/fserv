CREATE TABLE users (
	username VARCHAR(255) PRIMARY KEY,
	password VARCHAR(255),
	last_login TIMESTAMP
);

CREATE TABLE files (
	url_id VARCHAR(255) PRIMARY KEY,
	file_path VARCHAR(255),
	file_name VARCHAR(255),
	owner VARCHAR(255) REFERENCES users (username),
	expiration_date TIMESTAMP
);
