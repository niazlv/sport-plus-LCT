# Sport+

README на [English](/README.md)

Это бэкенд приложения SPORT+ для хакатона LCT 2024. Проект находится на стадии разработки.

## Обзор

Бэкенд Sport+ предназначен для поддержки тренеров в создании и управлении своими курсами, к которым пользователи могут получить доступ и участвовать в них. Приложение предоставляет несколько функций, включая управление пользователями, создание курсов, планирование занятий, функциональность чата и многое другое. Оно построено с акцентом на чистую архитектуру и использует современные технологии для бесшовной интеграции и масштабируемости.

## Основная разработка

Основная разработка этого проекта ведется на GitLab. Вы можете найти последние обновления и внести свой вклад в проект, посетив наш репозиторий на GitLab:

[Sport Plus LCT на GitLab](https://gitlab.sorewa.ru:12345/niaz/sport-plus-LCT)

## Как запустить приложение

Клонируйте репозиторий:

```bash
git clone https://github.com/niazlv/sport-plus-LCT.git
```

Перейдите в директорию проекта:

```bash
cd sport-plus-LCT
```

Создайте конфигурационный файл в каталоге `env`, например, создайте `dev.env`:

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

Или скопируйте и переименуйте `dev.example.env` и отредактируйте в соответствии с вашими требованиями:

```bash
mkdir env
cd env
cp dev.example.env dev.env
vim dev.env
```

Запустите файл Docker Compose:

```bash
COMPOSE_PROJECT_NAME=sport-plus docker-compose -f docker/dev.docker-compose.yml up -d --build
```

Или просто запустите:

```bash
make start-dev
```

Теперь вы можете получить доступ к бэкенду по адресу: [http://localhost:8080/v1](http://localhost:8080/v1)

## Документация и URL

- [openAPI.json](http://sport-plus.sorewa.ru:8080/openapi.json)
- [Swagger UI](http://sport-plus.sorewa.ru:8080/swagger)
- [HASURA panel](http://sport-plus.sorewa.ru:8085)

## Используемые технологии

- **Язык программирования:** Go
- **Веб-фреймворк:** Gin
- **Документация API:** Fizz, Swagger, Tonic
- **Аутентификация:** JWT
- **ORM:** GORM
- **База данных:** PostgreSQL
- **Контейнеризация:** Docker, Docker Compose
- **CI/CD:** GitLab CI/CD

## Возможности

- **Управление пользователями:** Позволяет тренерам и пользователям регистрироваться, входить в систему и "управлять" своими профилями.
- **Управление курсами:** Тренеры могут создавать и управлять курсами, которые включают классы, уроки и упражнения.
- **Функциональность чата:** Пользователи могут участвовать в чатах, связанных с их курсами.
- **Интеграция с календарем:** Пользователи могут планировать и просматривать свои тренировки.
- **Загрузка файлов:** Поддерживает загрузку изображений и других файлов.
- **Аутентификация:** Использует JWT токены для безопасной аутентификации и авторизации.
- **Реальное время:** Использует WebRTC для реализации функций в реальном времени.

## Статус проекта

Проект находится на стадии разработки и предназначен для хакатона LCT 2024.

## СПИСОК ЗАДАЧ

- [x] auth/signup
- [x] auth/signin
- [x] JWT токены
- [ ] Реализовать кастомный логгер с ротацией логов
- [x] модель User
- [x] swagger
- [ ] Рефакторинг кода в пакеты (user, auth)
- [ ] Отделить JWT от auth
- [ ] Переключить JWT токены на RS256
- [x] Добавить /auth/onboarding? (database.User) PUT
- [ ] Исправить коды ошибок в fizz/tonic
- [x] Завершить README.md с подробной информацией о проекте и настройке
- [x] Перенести СПИСОК ЗАДАЧ в управление проектом
- [ ] Переписать Docker файлы
  - [ ] docker файл
  - [ ] docker compose
- [ ] Добавить LICENSE
- [ ] Добавить валидацию ввода для операций с базой данных (особенно для auth, user, getUserByID)
- [x] chat
	- [x] Переписать socket.io
	- [x] Удалить конечную точку для отправки сообщений(send message)
	- [x] Транслировать сообщения от пользователя всем (broadcast: chatid, messageid, message)
	- [x] Один сокет на всех, отправлять chatid при подключении пользователя
	- [x] Вложения как "string"
- [x] webrtc
- [x] курс->класс->уроки
- [ ] Решить циклическую зависимость между auth и chat и отделить chat от auth
- [ ] Создать чат для курсов (createchatfromcourse)
- [ ] Переключиться на альтернативную реализацию socket.io
- [x] Фильтрация упражнений
- [x] Разрешить несколько упражнений в уроках. Реализовать как []exercise
	- [x] Обновить метод PUT
- [ ] Исправить функциональность прогресса/статуса

## Ссылки

- [Сайт LCT](https://i.moscow/cabinet/lct/profile/my-teams)
- [Исходный код мобильного приложения Sport+](https://github.com/justmeowme/sport_app_lct)
