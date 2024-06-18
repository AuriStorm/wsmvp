# wsmvp
backend<=>>centrifugo<=>clients communication


## перед началом
положить `centrifugo_config.json` в корневую директорию
положить `local_backend.env` в корневую директорию


## билд & запуск
`docker compose -f compose.yml build && docker compose -f compose.yml up --force-recreate`


## контейнеры:
backend-app - питоновское приложение, отсылает в центрифуго по http пейлоад в каналы
backend-streams - го приложение из примера для демонстрации перехвата бидерект сообщения с клиента через стрим
centrifugo - непосредственно сама центрифуга
frontend - клиент на жсе (можно открывать в разных вкладках как нового юзера)


## как использовать
откыть http://localhost:3000/
откыть сетевые запросы
зайти в админку http://127.0.0.1:8000 (пароль в поле admin_secret в centrifugo_config.json)
ввести в поле для ввода id/uid произвольный текст латиницей/цифры
нажать 'ПОЛУЧИТЬ ТОКЕН' (должен вывести ответ бекенд ручки с токеном, в логах бека должна отобразиться информация по ручке /centrifugo/subscribe/)
нажать 'ЗАИНИТИТЬ ЦЕНТРИФУГО'
1. должно вывести из под какого юзера заиничено
2. в сетевых запросах должен появитсья вызов 101 сокет соединения
2. в логе backend-streams должно вывести bidirectional subscribe called и инфо о саб реквесте вторым сообщением
нажать 'ОТПРАВИТЬ ПРИВЕТ С КЛИЕНТА' (в логе backend-streams должно вывести data from client {"input":"hello from client"})
вызвать ручку `POST http://127.0.0.1:8081/centrifugo/send-hello/` (должно появитсья сообщение под Полученный пейлоад надписью)
вызвать ручку `POST http://127.0.0.1:8081/centrifugo/send/` с пейлоадом в теле `{"payload": {"qwe": 3}}`
также можно отсылать пейлоад в центрифуго админке или смотреть пресенс по каналу `space`


## полезные ссылки
про аутентификацию по JWT
https://centrifugal.dev/docs/server/authentication

питоновский апи до центрифуги
https://github.com/centrifugal/pycent

питоновский клиент центрифуги
https://github.com/centrifugal/centrifuge-python

проксирование сообщений на сервер
https://centrifugal.dev/docs/server/proxy

проксирование через стримы
https://centrifugal.dev/docs/server/proxy_streams

настройка реверс-прокси
https://centrifugal.dev/docs/tutorial/reverse_proxy
