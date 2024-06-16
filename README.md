# Sport+

README in [Russian language](/README-RU.md)

This is the SPORT+ app backend for the LCT hackathon 2024. The project is currently in the development stage.

## Overview

The Sport+ backend is designed to support trainers in creating and managing their courses, which users can then access and participate in. The application provides several features including user management, course creation, class scheduling, chat functionality, and more. It is built with a focus on clean architecture and utilizes modern technologies for seamless integration and scalability.

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
DB_HOST=${COMPOSE_PROJECT_NAME}_postgres-db-1 
DB_PORT=5432 

JWT_SECRET=my-super-secret-key
HASURA_GRAPHQL_DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB}
HASURA_GRAPHQL_ENABLE_CONSOLE=true
HASURA_GRAPHQL_DEV_MODE=true
HASURA_GRAPHQL_ADMIN_SECRET=12Qwerty123!
" > dev.env
cd ..
```

Or copy and rename `dev.example.env` and edit by your requirements:

```bash
mkdir env
cd env
cp dev.example.env dev.env
vim dev.env
```

Run the Docker Compose file:

```bash
COMPOSE_PROJECT_NAME=sport-plus docker-compose -f docker/dev.docker-compose.yml up -d --build
```

or just run:

```bash
make start-dev
```

Now you can access the backend at: [http://localhost:8080/v1](http://localhost:8080/v1)

## Documentation and URLs

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
- **CI/CD:** GitLab CI/CD

## Features

- **User Management:** Allows trainers and users to sign up, sign in, and manage their profiles.
- **Course Management:** Trainers can create and manage courses, which include classes, lessons, and exercises.
- **Chat Functionality:** Users can participate in chats related to their courses.
- **Calendar Integration:** Users can schedule and view their training sessions.
- **File Uploads:** Supports uploading of images and other files.
- **Authentication:** Utilizes JWT tokens for secure authentication and authorization.
- **Real-Time Communication:** Implements WebRTC for real-time communication features.

## Project Status

This project is currently in the development stage and is intended for the LCT hackathon 2024.

## TODO LIST

- [x] auth/signup
- [x] auth/signin
- [x] JWT tokens
- [ ] Implement custom Logger with log rotation
- [x] user
- [x] swagger
- [ ] Refactor code into packages (user, auth)
- [ ] Decouple JWT from auth
- [ ] Switch JWT tokens to RS256
- [x] Add /auth/onboarding? (database.User) PUT
- [ ] Fix error codes in fizz/tonic
- [x] Complete README.md with detailed project and setup information
- [x] Move TODO LIST to project management
- [ ] Rewrite Docker files
  - [ ] docker file
  - [ ] docker compose
- [ ] Add LICENSE
- [ ] Add input validation for database operations (particularly for auth, user, getUserByID)
- [x] chat
	- [x] Rewrite socket.io
	- [x] Remove endpoint for sending messages
	- [x] Broadcast messages from user to all (chatid, messageid, message)
	- [x] One socket, send chatid when user connects
	- [x] Attachments as "string"
- [x] webrtc
- [x] course->class->lessons
- [ ] Resolve cyclic dependency between auth and chat, and decouple chat from auth
- [ ] Create chat for courses (createchatfromcourse)
- [ ] Switch to an alternative socket.io implementation
- [x] Exercise filtering
- [x] Allow multiple exercises in lessons. Implement as []exercise
	- [x] Update PUT method
- [ ] Fix progress/status functionality

## Links

- [LCT site](https://i.moscow/cabinet/lct/profile/my-teams)
- [Sport+ mobile app source code](https://github.com/justmeowme/sport_app_lct)
