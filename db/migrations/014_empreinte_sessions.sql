DELETE FROM sessions;

ALTER TABLE sessions RENAME COLUMN token TO empreinte_token;
