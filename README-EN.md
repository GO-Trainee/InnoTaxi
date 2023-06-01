# InnoTaxi

All repositories, that are going to be created during internship project must be private.
They become a part of NDA, signed by you (Forbidden do make those repos public, ever)!

InnoTaxi - app for ordering a taxi.

In app realization there must be:

3 user roles:
1) User;
2) Driver;
3) Analyst.

2 wallet types:
Personal
Family

3 taxi types:
1) Economy;
2) Comfort;
3) Business.

App must be based on microservice architecture, with clean architecture principles.

4 microservices must be created:

1) User Service;
2) Driver Service;
3) Order Service;
4) Analytic Service.
   Detailed description of services provided below.

Functional requirements:

User Service:
1) User can sign up.
    - Fields for signing up: name, phone number, email, password.
2) User can sign in.
    - Fields for signing in: phone number, password.
3) User can log out from app (token, that was given to user must become unacceptable by system, even if it still valid).
4) User can order a taxi.
    - Fields for ordering a taxi: taxi type, from, to.
    - While ordering a taxi, system is seeking for a free drive(status = free). When drive if found, order will be created with in progress status. Drive becomes busy (status = busy).
    - If there were no available drivers, user must stay waiting for a free driver for some configurable amount of time. If after this time no driver was found, the user should receive a rejection response.
    - If no drivers are available, a queue is formed of waiting users. The first user to order a taxi should be the first to receive a driver or a message that no driver has been found.
5) User can rate the trip. Only the last trip can be rated. Only the most recent trip can be rated if less than a certain (configurable) amount of time has elapsed since the trip. Rating from 1 to 5 inclusive.
   User can leave an optional comment for the trip.
6) User can view their trips (taxi type, driver, from, to).
7) User can view their profile (name, phone number, email, rating).
8) User can update their profile (name, phone number, email).
9) User can delete their profile (из бд не удаляется, помечается как удаленный).
10) User must have a wallet, also must be added a feature of creating a “family” wallet (wallet for many users), family wallet must be attached to user's personal wallet, user can add new members to family wallet by phone number.
11) If there ara many wallets available for the User, User can choose from which balance will be drained.
12) User can cash in their wallet. Family wallet can be cashed in only by its owner.
13) User can view their wallet's transactions. For family wallet only owner can see transactions.
14) User can't order a taxi if there not enough money in the chosen waller.

Driver Service:
1) Driver может зарегистрироваться.
    - Поля для регистрации: name, phone number, email, password, taxi type.
2) Driver может войти в систему.
    - Поля для входа: phone number, password.
3) Driver может сменить свой статус с busy на free, когда поездка закончена. При этом статус order'a меняется на finished.
4) Driver может оценить поездку. Оценить можно только последнюю поездку, если с нее прошло менее некоторого(настраиваемого) времени. Оценка от 1 до 5 включительно.
5) Driver может узнать свой рейтинг. Рейтинг формируется на основе последних 20 оценок от юзеров.
6) Driver может просмотреть свои поездки (taxi type, user, from to).
7) У Driver могут быть статусы busy и free.

Order Service:
1) Order Service выступает в качестве оркестратора для User и Driver сервисов для создания заказов.
2) Можно получить список заказов. По опциональным параметрам (может быть любое поле заказа).
3) Поля Order: user, driver, from, to, taxi type, date, status, comment.
4) У Order могут быть статусы in progress и finished.
5) Analyst имеет возможность поиска по списку заказов.
6) Order Service хранит в себе цены для каждого из taxi type.


Analytic Service:
1) Analyst может смотреть статистику заказов(количество по дням, месяцам).
2) Analyst может смотреть рейтинг всех водителей.
3) Analyst может смотреть рейтинг всех пользователей.
4) Analyst имеет заранее созданный в системе аккаунт, логин происходит по username, password.
5) Все регистрации и выполненные заказы должны записываться в аналитическую базу.



Схема приложения:

<img src="Design-diagram-v3.png" width="600" height="500" /> 

Допускается добавление новых сервисов, уменьшать количество сервисов не допускается. При изменении схемы(добавление новых сервисов) приложения необходимо дополнить диаграмму и выложить себе в репозиторий в .png и .drawio форматах.

Нефункциональные требования:

1) Общее:
- Работа по GitHub Flow.
- Две основных ветки (main, dev).
- Main выступает в качестве релизной ветки. Сюда должны заливаться только стабильные версии (с реализованным функционалом).
- При добавление новой фичи, рефакторинге или фиксинге багов необходимо отбренчиваться от Dev ветки и создавать новую.
- Делаем работу в созданных ветках и создаем Pull Request в Dev-ветку.
- Когда сделали pull request, пишем ментору, приступаем к новой задаче не дожидаясь ревью. Работаем по Jira флоу, описанному в Onboarding тикете.
- Если ментор оставит комментарии на pr, исправляем и заливаем изменения
- Если ментор аппрувит pr то мерджим его в dev. Перед мерджем в Dev необходимо сделать Squash коммитов.

2) User service:
- PostgreSQL в качестве БД.
- Таблицы Users, Trips, Wallets, Transactions.
- Wallets имеет MnM связь с Users. В Wallets есть колонка Type, которая является Enum. Может принимать значения Personal и Family.
- В PostgreSQL должны храниться последние 20 заказов для каждого пользователя. При превышении количества > 20 записей неактуальные чистятся триггером.
- На основе 20 оценок с последних заказов у пользователя формируется средний рейтинг.
- Таблица со всеми транзакциями по списанию баланса. Статусы транзакции: create, blocked, success and canceled. При создании пользователем заказа ему присваивается статус **create**, если баланс пользователя >= стоимости поездки статус меняется на **blocked**, в ином случае на **canceled**. По завершению поездки статус переводится в **success**.
- Redis в качестве кеша. Реализовать хранение токена пользователя и logout (удаление токена).
- Prometheus и Grafana для сбора метрик.
- Swagger генерируется из кода приложения.
- Front-end: Использовать Vue.js (3.0) composition API. Создать форму регистрации, авторизации, просмотр профиля, изменение полей в профиле, удаление профиля. Использовать компоненты, Pinia.
- VCS: GitHub; CI/CD: Github Actions.

3) Driver service:
- Код транспортного слоя генерируется из Swagger описания. (Не пишется руками HTTP взаимодействие).
- MongoDB в качестве БД.
- В базе хранится: информация о Driver (name, phone number, email, password, taxi type), информация о последних 20 заказах, рейтинг водителя в виде массива, баланс водителя (заработок). Когда длина массива с заказами становится >20 необходимо удалять неактуальные.
- Агрегация рейтинга драйвера на основе 20 последних оценок.
- Должны быть добавлены pprof хендлеры (и использованы). Или другой тип профилирования, если используется не GO.
- Front-end: Использовать Angular (latest), создать форму регистрации, авторизации и просмотр профиля, изменение полей в профиле, удаление профиля.
- VCS: Gitlab; CI/CD: Gitlab CI/CD;

4) Order service:
- GraphQL в качестве транспортного слоя.
- Эндпоинт для поиска по полям: DriverID, UserID, from, to, fromDate, toDate. Реализовать частичный поиск (когда некоторые поля не указаны - должны игнорироваться).
- Добавить пагинацию для данного эндпоинта. Поля offset; limit.
- Elasticsearch для order service (возможность поиска по from или to созданных поездок). Поиск по префиксу, full text search по комментариям к поездкам, добавить возможность поиска с транслитерацией и лексическими ошибками.
- Для реализации ожидания юзером свободных водителей использовать многопоточное программирование.
- Доп задание: реализовать возможность подбора водителей исходя из рейтинга пользователя.
- Front-end: Использовать React + redux/redux toolkit. Создать main page (фронт к GQL эндпоинтам) с фильтрами (from - to; fromDate - toDate; User name; Driver name; taxi type etc.).
- VCS: BitBucket; CI/CD: Bitbucket pipelines.

5) Analyst service:
- ClickHouse в качестве БД.
- Читает из кафки события и пишет в ClickHouse.
- VCS: Github; CI/CD: Circle CI.

Technical requirements:
- Kafka в качестве message broker.
- Для RPC взаимодействия использовать GRPC.
- Для каждого сервиса должен быть описан Dockerfile.
- Поднятие всего приложения осуществлять через docker-compose.
- Все изменяемые переменные (подключение к бд, время ожидания юзером свободного водителя) должны устанавливаться через environment variables.
- Каждый сервис должен иметь README (описано, что необходимо для запуска, environment variables) и swagger (подробно описаны все эндпоинты).
- Все сервисы должны быть покрыты интеграционными тестами.
- Для каждого сервиса должен быть настроен CI/CD со степами: 1) тесты, 2) линтер, 3) линтер для протофайлов, 4) vulncheck, 5) билд образа и залив на свой dockerHub (в мастер ветке).
- Аутентификация с использованием JWT.

  Go:
- UserService: Gin в качетсве HTTP-библиотеки. Analyst: Fiber.
- Использовать golangci-lint
- Для работы с Postgres допустимо использовать чистые запросы или sqlc, sqlx. Использование ORM не допускается.
- Для реализации ожидания юзером свободных водителей использовать горутины и каналы, при необходимости пакеты sync, x/sync.
- Для миграций использовать Goose или go-migrate.
- Style guidlines [сюда](https://rakyll.org/style-packages/)
- Для тестирования использовать табличные тесты, gomock, testify/suite, ginkgo, gomega.
