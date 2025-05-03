# PAPS-lab

В данном репозитории представлен веб-сервер, реализованный на языке Go в рамках курса "Проектирование и архитектура программных систем" на тему "Телефонный справочник".

## Возможности

- удобный графический интерфейс для работы с приложением
- хранение сущностей "абонент" в базе данных `Postgress`
- менеджер сессиий с сохранением сессий в `Redis`
- менеджер паролей
- инструменты для нагрузочного тестирования с помощью Prometheus/Grafana, находящиеся в ветке `testing`

## Требования

- Go версии 1.16 и выше
- Docker и Docker-compose

## Описание модулей

1. `basic_types` - базовые типы для работы с базой данных и Redis;
2. `studiodb` - подключение к базе данных Postgress;
3. `password_manager` - менеджер паролей, взаимодействующий с базой данных;
4. `register` - регистр телефонного справочника, взаимодействующий с базой данных;
5. `session_manager` - менеджер сессий, взаимодействующий с Redis.

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
   
4. Запустите контейнеры через Docker-compose
   ```bash
   docker-compose up --build -d
   ```

5. Установите зависимости
   ```bash
   go mod download
   ```

6. Запустите веб-сервер
   ```bash
   go run .
   ```

   После успешного запуска в консоль будет выведено сообщение `Starting server at :8080`.
   
После запуска откройте `localhost:8080`
