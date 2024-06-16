start-dev:
	COMPOSE_PROJECT_NAME=sport-plus docker-compose -f docker/dev.docker-compose.yml up -d --build
