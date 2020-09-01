# Avoid pg_dump version mismatch where multiple postgres versions are
# installed by specifying the absolute path.
PG_DUMP=/usr/local/opt/postgresql@12/bin/pg_dump
SCHEMA=./db/schema.sql
db_dump_schema:
	@if [ $(FB05_ENV) == "development" ]; then\
		echo "Dumping schema...";\
        $(PG_DUMP) fb05_$(FB05_ENV) --file=$(SCHEMA) --schema-only;\
		echo "Schema dumped to $(SCHEMA)\n";\
	else\
		echo "Schema should only be dumped in development.";\
    fi

db_migrate_down:
	go run cmd/migrate/migrate.go down
	@$(MAKE) db_dump_schema

db_drop:
	go run cmd/migrate/migrate.go drop

db_force_version:
	go run cmd/migrate/migrate.go force $(VERSION)

db_show_version:
	go run cmd/migrate/migrate.go version

db_migrate_steps:
	go run cmd/migrate/migrate.go steps $(STEPS)
	@$(MAKE) db_dump_schema

db_migrate_to_version:
	go run cmd/migrate/migrate.go toVersion $(VERSION)
	@$(MAKE) db_dump_schema

db_migrate_up:
	go run cmd/migrate/migrate.go up
	@$(MAKE) db_dump_schema

db_test_prepare:
	psql --set ON_ERROR_STOP=on fb05_test < $(SCHEMA)

run_dev:
	go run cmd/fb05/main.go