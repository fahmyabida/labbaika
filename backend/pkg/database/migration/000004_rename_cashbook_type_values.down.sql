ALTER TABLE cashbooks DROP CONSTRAINT cashbooks_type_check;

UPDATE cashbooks SET type = 'CREDIT' WHERE type = 'IN';
UPDATE cashbooks SET type = 'DEBIT' WHERE type = 'OUT';

ALTER TABLE cashbooks ADD CONSTRAINT cashbooks_type_check CHECK (type IN ('CREDIT', 'DEBIT'));
