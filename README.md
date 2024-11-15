# Менеджер паролей GophKeeper

GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю надёжно и безопасно хранить логины, пароли, бинарные данные и прочую приватную информацию.

Сервер поддерживает следующий функционал:

* регистрация, аутентификация и авторизация пользователей;
* хранение приватных данных пользователей;
* синхронизация данных между несколькими авторизованными клиентами одного владельца;
* передача приватных данных владельцу по запросу.

Клиент реализует следующую бизнес-логику:

* регистрация, аутентификация и авторизация пользователей на удалённом сервере;
* доступ к приватным данным по запросу.

## Настройка и запуск сервера и клиента

Перед запуском сервера необходимо задать переменные окружения:

Сервер:
```
ADDRESS=127.0.0.1:9090
DATABASE_URI=postgres://postgres:postgres@localhost:5432/postgres
AUTH_KEY=xiuw1bi4r98vd1(&*6
AUTH_EXPIRATION_TIME=25
```

Клиент:
```
ADDRESS=127.0.0.1:9090
ENCRYPTION_KEY=WYJcWgkItShq513L21E1CFuz6uQWDy3p
```

## Процедуры регистрации, аутентификации, авторизации

При регистрации пользователя необходимо указать логин и пароль.
Пример команды регистрации:

```
./gophkeeper auth register --username=user@mail.com --password=111
```

В случае успешного выполнения запроса регистрации нового пользователя, сервер вернет в ответ токен доступа. 
Токен будет сохранен в текстовый файл.

В случае необходимости, токен доступа можно запросить повторно с помощью команды:

```
./gophkeeper auth login --username=user@mail.com --password=111
```

## Хранение приватных данных пользователя

Для каждой операции необходимо иметь актуальный токен

### Добавление приватных данных

1. Пример команды, сохраняющей данные банковской карты:

```
./gophkeeper secret create card --name=MyFirstCard  \
    --number=1234 5678 9012 3456 \
    --date=12/24 \
    --code=123 \
    --holder=Alex B.
```

2. Пример команды создания пары логин-пароль:

```
./gophkeeper secret create credentials \
  --name=yandex-practicum \
  --login=user@mail.com \
  --password=12345678
```

3. Пример команды для сохранения текстовой информации

```
./gophkeeper secret create text \
  --name=Small text \
  --data="Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua"
```


### Получение данных

Пример команды получения секрета по имени:

```
./gophkeeper secret read --name="Small text"
```

### Удаление данных

Пример команды удаления данных:

```
./gophkeeper secret delete --name=MyFirstCreds
```