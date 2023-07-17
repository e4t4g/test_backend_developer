init_db: ## Initializing DB
	(docker stop postgres && docker rm postgres) || true
	docker run \
	--name postgres \
	-e POSTGRESQL_USERNAME=postgres \
	-e POSTGRESQL_DATABASE=spots \
	-e POSTGRESQL_PASSWORD=123 \
	-p 5432:5432 \
	bitnami/postgresql:latest

fill_db: ## Adding extension, fill the table with data
	docker exec -e PGPASSWORD=123 -it postgres psql -U postgres -d spots -c "CREATE EXTENSION postgis;" ; \
	docker exec -e PGPASSWORD=123 -i postgres psql -U postgres -d spots < storage/spots.sql

output: ## Return 3 columns where a count of domain is greater than 1, website field only contains domain, it shows how many spots have the same domain
	docker exec -e PGPASSWORD=123 -it postgres psql -U postgres -d spots -c "COPY (SELECT LOWER(name) AS \"spot name\", CASE WHEN website LIKE 'www.%' THEN REPLACE(SPLIT_PART(website, 'www.', 2), '/', '') WHEN website LIKE '%://%' THEN REPLACE(SPLIT_PART(SPLIT_PART(website, '://', 2), '/', 1), 'www.', '') ELSE replace(split_part(website, '/', 1), '/', '') END AS domain, COUNT(*) AS \"count number for domain\" FROM \""MY_TABLE"\" GROUP BY LOWER(name), domain HAVING COUNT(website) > 1) TO '/docker-entrypoint-initdb.d/output_1task.csv' CSV HEADER"
	docker cp postgres:/docker-entrypoint-initdb.d/output_1task.csv result/