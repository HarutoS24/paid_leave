.PHONY: up down restart_backend restart reset backend_restart backend frontend db
up:
	docker compose up -d
down:
	docker compose down
restart:
	docker compose restart
reset:
	docker compose down -v --rmi all
backend_restart:
	docker compose restart backend
backend:
	docker compose exec -it backend /bin/bash
frontend:
	docker compose exec -it frontend /bin/bash
db:
	docker compose exec -it db /bin/bash