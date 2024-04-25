ALTER TABLE visits
    DROP CONSTRAINT visits_link_id_fkey;

ALTER TABLE visits
    ADD CONSTRAINT visits_link_id_fkey FOREIGN KEY (link_id) REFERENCES links (id) ON DELETE CASCADE;
