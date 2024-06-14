import requests
import json
from datetime import datetime, timedelta

# URL для авторизации и создания расписания
BASE_URL = "http://sport-plus.sorewa.ru:8080/v1"
AUTH_URL = f"{BASE_URL}/auth/signin"
CALENDAR_URL = f"{BASE_URL}/calendar"

# Данные пользователей
users = [
    {"login": "coach_acc0", "password": "12345", "role": "coach", "id":41},
    {"login": "client_calendar", "password": "client_calendar", "role": "client", "id":40}
]

# Функция для авторизации и получения токена
def get_token(login, password):
    response = requests.get(AUTH_URL, params={"login": login, "password": password})
    if response.status_code == 200:
        return response.json().get("token")
    else:
        raise Exception(f"Failed to authenticate {login}: {response.text}")

# Функция для создания мероприятия
def create_schedule(token, data):
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.post(CALENDAR_URL, headers=headers, json=data)
    if response.status_code not in [200, 201]:
        raise Exception(f"Failed to create schedule: {response.text}")

# Генерация мероприятий
def generate_schedules():
    now = datetime.utcnow()  # Используем UTC время
    schedules = []

    # Локальные мероприятия
    local_events = [
        {"days": 1, "title": "Завтрашняя тренировка"},
        {"days": 3, "title": "Тренировка через 3 дня"},
        {"days": 7, "title": "Тренировка через неделю"},
        {"days": 30, "title": "Тренировка через месяц"}
    ]

    for user in users:
        for event in local_events:
            for _ in range(2):  # Два мероприятия на каждого пользователя
                start_time = now + timedelta(days=event["days"], hours=12)
                end_time = start_time + timedelta(hours=1)
                schedules.append({
                    "client_id": user["id"],
                    "date": start_time.isoformat(timespec='seconds') + 'Z',
                    "start_time": start_time.isoformat(timespec='seconds') + 'Z',
                    "end_time": end_time.isoformat(timespec='seconds') + 'Z',
                    "type": "local",
                    "reminder_client": True,
                    "reminder_coach": True,
                    "is_global": False,
                })

    # Глобальные мероприятия
    for _ in range(5):
        start_time = now + timedelta(days=7, hours=12)
        end_time = start_time + timedelta(hours=1)
        schedules.append({
            "client_id": 41,
            "date": start_time.isoformat(timespec='seconds') + 'Z',
            "start_time": start_time.isoformat(timespec='seconds') + 'Z',
            "end_time": end_time.isoformat(timespec='seconds') + 'Z',
            "type": "global",
            "reminder_client": True,
            "reminder_coach": True,
            "is_global": True,
        })

    return schedules

def main():
    try:
        # Авторизация и получение токенов
        tokens = {user["login"]: get_token(user["login"], user["password"]) for user in users}

        # Генерация мероприятий
        schedules = generate_schedules()

        # Создание мероприятий
        for schedule in schedules:
            if schedule["is_global"]:
                token = tokens["coach_acc0"]
            else:
                token = tokens["client_calendar"] if schedule["client_id"] == 1 else tokens["coach_acc0"]
            create_schedule(token, schedule)

        print("All schedules created successfully!")

    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    main()