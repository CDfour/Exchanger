Документация
------------

Приложение Exchanger раз в 15 минут отправляет запрос на биржу https://exchangeratesapi.io/ для получения актуального курса Рубля(RUB), Евро(EUR) и Юаня(CNY).
Для использования нужно получить свой apikey.
Полученные результаты записываются в базу данных.
Сервер обрабатывает два эндпоинта, по одному выводит информация о валютах которые имеются в базе данных, по второму выводит актуальные курсы для указанных валют в параметре URL.
Также можно указать время на которое нужен курс в формате гггг-мм-дд чч-мм-сс
http://localhost:8080/currencies
http://localhost:8080/rates?symbols=RUB&time=2022-09-20%2020:00:00
Для создания сервера и обработки запросов используется пакет gin. https://github.com/gin-gonic/gin
Для запуска планировщика используется пакет cron. https://github.com/robfig/cron
В качестве драйвера PostgreSQL выбор пал на pgx. https://github.com/jackc/pgx
Для логирования использовался пакет logrus. "github.com/sirupsen/logrus"
Документация реализована с помощью swagger. https://github.com/swaggo/swag

Контактная почта - biv_1998@mail.ru