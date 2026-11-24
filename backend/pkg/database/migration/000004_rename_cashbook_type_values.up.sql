ALTER TABLE cashbooks DROP CONSTRAINT cashbooks_type_check;

UPDATE cashbooks SET type = 'IN' WHERE type = 'CREDIT';
UPDATE cashbooks SET type = 'OUT' WHERE type = 'DEBIT';

ALTER TABLE cashbooks ADD CONSTRAINT cashbooks_type_check CHECK (type IN ('IN', 'OUT'));
