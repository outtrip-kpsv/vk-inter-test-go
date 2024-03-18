приложение разворачивается по команде `make run` \
(*переменные окружения подключаемые к контейнерам находятся в деректории* `env`)\
фильмотека разделена на слои:
- транспорт `internal\io`
- репозиторий `internal\db`
- бизнес логика `internal\bl`

_в репозитории находятся файлы миграции для создания таблиц и первоначального заполнения_\
(создан пользователь `admin` с паролем `admin` с ролью админ для создания и изменения записей)\
остальные пользователи создаются (см документацию) с ролью юзер

все роуты описаны в `internal\io\routes.go` к авторизации и создании пользователя есть доступ у всех\
в результате отработки этих хендлеров можно получит bearer токен который нужно передавать для доступа к остальным хендлерам\
хендлеры в которых происходит изменение бд дополнительно проверяют пользователя на наличие необходимой роли "admin"

в контейнере поднимается:
- база postgresql,
- приложение с сервером, (на порту 3000)
- swagger документация (code-first) по API доступная по адресу `http://localhost:8085/`

поктыто тестами слой с "бизнес логикой" и пакет с утилитами, команда `make test`