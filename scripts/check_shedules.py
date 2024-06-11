import requests

# URL для авторизации и получения расписаний
BASE_URL = "http://sport-plus.sorewa.ru:8080/v1"
AUTH_URL = f"{BASE_URL}/auth/signin"
CALENDAR_URL = f"{BASE_URL}/calendar"

# Данные пользователей
users = [
    {"login": "coach_acc0", "password": "12345", "role": "coach"},
    {"login": "client_calendar", "password": "client_calendar", "role": "client"}
]

# Функция для авторизации и получения токена
def get_token(login, password):
    response = requests.get(AUTH_URL, params={"login": login, "password": password})
    if response.status_code == 200:
        return response.json().get("token")
    else:
        raise Exception(f"Failed to authenticate {login}: {response.text}")

# Функция для получения расписания
def get_schedules(token, endpoint):
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(f"{CALENDAR_URL}/{endpoint}", headers=headers)
    if response.status_code == 200:
        return response.json()
    else:
        raise Exception(f"Failed to get schedules: {response.text}")

def main():
    try:
        # Авторизация и получение токенов
        tokens = {user["login"]: get_token(user["login"], user["password"]) for user in users}

        # Получение глобальных расписаний
        global_schedules = get_schedules(tokens["coach_acc0"], "global")
        print("Global Schedules:")
        for schedule in global_schedules:
            print(schedule)

        # Получение локальных расписаний для тренера
        coach_local_schedules = get_schedules(tokens["coach_acc0"], "local")
        print("\nCoach Local Schedules:")
        for schedule in coach_local_schedules:
            print(schedule)

        # Получение локальных расписаний для клиента
        client_local_schedules = get_schedules(tokens["client_calendar"], "local")
        print("\nClient Local Schedules:")
        for schedule in client_local_schedules:
            print(schedule)

    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    main()
