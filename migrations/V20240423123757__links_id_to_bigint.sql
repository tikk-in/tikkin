ALTER TABLE links
    ALTER COLUMN id TYPE BIGINT;

ALTER SEQUENCE links_id_seq AS BIGINT MAXVALUE 9223372036854775807;