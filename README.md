# Sport+

This is SPORT+ app backend for LCT hackaton 2024.

## Main Development

The main development of this project is conducted on GitLab. You can find the latest updates and contribute to the project by visiting our GitLab repository:

[Sport Plus LCT on GitLab](https://gitlab.sorewa.ru:12345/niaz/sport-plus-LCT)

## How to run app

Clone repo:

```bash
git clone https://github.com/niazlv/sport-plus-LCT.git
```

Go into it:

```bash
cd sport-plus-LCT.git
```

Create configure file into env, for example create dev.env

```bash
mkdir env
cd env
echo "POSTGRES_USER=postgres 
POSTGRES_PASSWORD=postgres 
POSTGRES_DB=postgres 

DB_HOST=docker-db-1 
DB_PORT=5432 

JWT_SECRET=my-super-secret-key
" > dev.env
cd ..
```

And Run docker compose file.

```bash
docker-compose -f docker/dev.docker-compose.yml up -d --build
```

Now you can access to Backend by: http://localhost:8080/v1

## Documentation

- [openAPI.json](http://sport-plus.sorewa.ru:8080/openapi.json)
- [swagger](http://sport-plus.sorewa.ru:8080/swagger)

## TODO LIST

- [x] auth/signup
- [x] auth/signin
- [ ] JWT tokens
- [x] user
- [x] swagger
- [ ] вынести участки кода в пакеты(user,auth)
- [ ] JWT вынести из(auth)
- [ ] JWT tokens переделать под RS256
- [ ] добавить /auth/onboarding? (database.User) PUT
- [ ] fizz/tonic. Починить коды ошибок.
- [x] README.md. Дописать открытую документацию на запуск и о проекте, что использую.
- [x] Перенести TODO LIST в проект
- [ ] Переписать docker
  - [ ] docker file
  - [ ] docker compose

## links

- [lct site](https://i.moscow/cabinet/lct/profile/my-teams)
- [Sport+ mobile app src](https://github.com/justmeowme/sport_app_lct)
