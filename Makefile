.PHONY: start
start:
	mv -n backend/config/config.example.toml backend/config/config.toml
	docker-compose up --build
