-- name: InsertVisit :one
INSERT INTO visits (id, link_id, user_agent, referrer, country_code)
VALUES (@id, @link_id, @user_agent, @referrer, @country_code)
RETURNING *;

-- name: CountVisitsByLinkID :one
SELECT COUNT(*)
FROM visits
WHERE link_id = @link_id;
