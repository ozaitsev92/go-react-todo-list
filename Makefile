.PHONY: start
start: .env ## Build and run the full stack with docker compose
	docker compose up --build

.PHONY: env
env: .env

# Local secrets, generated once; delete .env to rotate.
.env:
	@umask 077; { \
	  echo "JWT_SIGNING_KEY=$$(openssl rand -hex 32)"; \
	  echo "MONGO_ROOT_USER=todo"; \
	  echo "MONGO_ROOT_PASSWORD=$$(openssl rand -hex 24)"; \
	} > .env
	@echo "generated .env with fresh secrets"
