# Sport+

This is the SPORT+ app backend for the LCT hackathon 2024. The project is currently in the development stage.

## Main Development

The main development of this project is conducted on GitLab. You can find the latest updates and contribute to the project by visiting our GitLab repository:

[Sport Plus LCT on GitLab](https://gitlab.sorewa.ru:12345/niaz/sport-plus-LCT)

## How to Run the App

Clone the repository:

```bash
git clone https://github.com/niazlv/sport-plus-LCT.git
```

Navigate into the project directory:

```bash
cd sport-plus-LCT
```

Create a configuration file in the `env` directory, for example, create `dev.env`:

```bash
mkdir env
cd env
echo "POSTGRES_USER=postgres 
POSTGRES_PASSWORD=postgres 
POSTGRES_DB=postgres 

COMPOSE_PROJECT_NAME=sport-plus
DB_HOST=docker-db-1 
DB_PORT=5432 

JWT_SECRET=my-super-secret-key


HASURA_GRAPHQL_DATABASE_URL:postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB}
HASURA_GRAPHQL_ENABLE_CONSOLE:true
HASURA_GRAPHQL_DEV_MODE:true
HASURA_GRAPHQL_ADMIN_SECRET:12Qwerty123!
" > dev.env
cd ..
```

Run the Docker Compose file:

```bash
docker-compose -f docker/dev.docker-compose.yml up -d --build
```

Now you can access the backend at: [http://localhost:8080/v1](http://localhost:8080/v1)

## Documentation and urls

- [openAPI.json](http://sport-plus.sorewa.ru:8080/openapi.json)
- [Swagger UI](http://sport-plus.sorewa.ru:8080/swagger)
- [HASURA panel](http://sport-plus.sorewa.ru:8085)

## Technologies Used

- **Programming Language:** Go
- **Web Framework:** Gin
- **API Documentation:** Fizz, Swagger, Tonic
- **Authentication:** JWT
- **ORM:** GORM
- **Database:** PostgreSQL
- **Containerization:** Docker, Docker Compose

## Project Status

This project is currently in the development stage and is intended for the LCT hackathon 2024.

## TODO LIST

- [x] auth/signup
- [x] auth/signin
- [ ] JWT tokens
- [x] user
- [x] swagger
- [ ] вынести участки кода в пакеты(user,auth)
- [ ] JWT вынести из(auth)
- [ ] JWT tokens переделать под RS256
- [x] добавить /auth/onboarding? (database.User) PUT
- [ ] fizz/tonic. Починить коды ошибок.
- [x] README.md. Дописать открытую документацию на запуск и о проекте, что использую.
- [x] Перенести TODO LIST в проект
- [ ] Переписать docker
  - [ ] docker file
  - [ ] docker compose
- [ ] Дописать LICENCE
- [ ] Добавить защит от дурака, при обращении к БД(в частности для auth, user, getUserByID)
- [x] chat
	- [x] rewrite socket.io
	- [x] delete endpont send messages
	- [x] broadcast message from user to all (chatid, messageid, message)
	- [x] one socket, send chatid when user connect.
	- [x] attachments "string"
- [x] webrtc
- [x] course->class->lessions
- [ ] решить циклическую зависимость auth<->chat и вынести chat из auth
- [ ] TODO: createchatfromcourse. Чат для курса
- [ ] перейти на другой socket.io
- [x] фильтрация exercise

## Links

- [LCT site](https://i.moscow/cabinet/lct/profile/my-teams)
- [Sport+ mobile app source code](https://github.com/justmeowme/sport_app_lct)
