# PAPS-lab

В данном репозитории представлен веб-сервер, реализованный на языке Go в рамках курса "Проектирование и архитектура программных систем" на тему "Телефонный справочник".

## Описание модулей

1. `handler` - обработчики http запросов;
2. `studiodb` - работа с базой данных PostgreSQL;
3. `session_manager` - менеджер сессий, взаимодействующий с Redis.

## Установка и запуск

1. Клонируйте репозиторий

   ```bash
   git clone https://github.com/imightbuyaboat/PAPS-lab
   cd PAPS-lab
   ```
   
2. В корне проекта создайте `.env` файл

   ```bash
   nano .env
   ```

   со следующим содержимым:

   ```env
   SQL_HOST=localhost
   SQL_PORT=5432
   SQL_DB=your_data_base
   SQL_USER=your_user
   SQL_PASSWORD=your_password

   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=your_password
   ```
   
3. Запустите сервер
   ```bash
   docker compose -f docker-compose.yml --env-file .env up --build -d
   ```

   После успешного запуска в консоль будет выведено сообщение `Starting server at :8080`.
   
После запуска откройте `localhost:8080`
