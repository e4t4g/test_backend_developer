SELECT
       LOWER(name) AS "spot name",
       CASE WHEN website LIKE 'www.%'
             THEN REPLACE(SPLIT_PART(website, 'www.', 2), '/', '')
           WHEN website LIKE '%://%'
               THEN REPLACE(SPLIT_PART(SPLIT_PART(website, '://', 2), '/', 1), 'www.', '')
           ELSE replace(split_part(website, '/', 1), '/', '') END AS domain,
       COUNT(*) AS "count number for domain"
        FROM "MY_TABLE"
        GROUP BY LOWER(name), domain
        HAVING COUNT(website) > 1