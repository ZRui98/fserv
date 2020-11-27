ALTER TABLE files ADD is_private BOOLEAN;
UPDATE files SET is_private = false;