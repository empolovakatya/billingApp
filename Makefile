build:
	docker-compose build billingapp

run:
	docker-compose up billingapp

createdb:
	docker exec -it billingapp_db_1 createdb --username=postgres --owner=postgres billing_db

dropdb:
	docker exec -it billingapp_db_1 dropdb --username=postgres billing_db

migrate:
	migrate -path ./schema -database 'postgresql://postgres:v&487fnd4jbvf8@0.0.0.0:5436/billing_db?sslmode=disable' up
