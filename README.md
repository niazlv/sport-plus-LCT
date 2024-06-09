## How to run app
```bash
docker-compose -f docker/dev.docker-compose.yml up -d --build
```

## Documentation
- [openAPI](http://sport-plus.sorewa.ru:8080/openapi.json)
- [swagger](http://sport-plus.sorewa.ru:8080/swagger)

## TODO LIST
- [x] auth/signup
- [x] auth/signin
- [x] JWT tokens
- [x] user
- [x] swagger
- [ ] вынести участки кода в пакеты(user,auth)
- [ ] JWT вынести из(auth)
- [ ] JWT tokens переделать под RS256
- [ ] добавить /auth/onboarding? (database.User) PUT
- [ ] fizz/tonic. Починить коды ошибок.
- [ ] README.md. Дописать открытую документацию на запуск и о проекте, что использую.
- [x] Перенести TODO LIST в проект
