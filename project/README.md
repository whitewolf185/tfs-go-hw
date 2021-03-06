# Инструкции по использованию курсового проекта

------

## Начало работы
Для начала нужно подготовить несколько ENV переменных:
* TOKEN_PATH_PRIVATE -- записывается путь к файлу, где лежит private api токен
* TOKEN_PATH_PUBLIC -- записывается путь к файлу, где лежит public api токен
* WS_URL -- URL веб сокета, откуда будут тянуться свечки
* TG_BOT_TOKEN -- путь к файлу, где лежит токен телеграм бота

Например
```text
TOKEN_PATH_PRIVATE=D:/GO/tfs-go-hw/project/non-project-files/testing_private.txt;
TOKEN_PATH_PUBLIC=D:/GO/tfs-go-hw/project/non-project-files/testing_public.txt;
WS_URL=wss://demo-futures.kraken.com/ws/v1?chart
TG_BOT_TOKEN=D:/GO/tfs-go-hw/project/non-project-files/testing_tg.txt
```

Для запуска курсового проекта, нужно запустить файл `main.go` (находится в корне папки проекта)

## Взаимодействие
### Начальная настройка
Все взаимодействие с приложением происходит через телеграм [бота](http://t.me/Tfs_CP_bot)

Для начала нужно запустить ТГ бота и ввести команду `/start`.

Бот не будет работать, если не ввести ему начальные настройки. Для этого нужно ввести команду `/option`. 
Далее бот попросит 2 параметра:
1. Тикет 
2. Период свечки

**Список поддерживаемых тикетов:**
* `PI_XBTUSD`
* `PI_ETHUSD`
Вводить нужно с точностью до регистра

**Список поддерживаемых периодов:**
* `1m`
* `2m`
* `1h`

### Использование
Бот поддерживает несколько функций:
1. Покупка одного тикета, указанного ранее в настройках по маркет цене. Исполняется командой `/buy`
2. Продажа одного тикета, указанного ранее в настройках по маркет цене. Исполняется командой `/sell`
3. Настройка индикатора

**Поддерживаемые индикаторы:**
* _stoploss_. Исполняется командой `/stoploss`
* _takeprofit_. Исполняется командой `/takeprofit`

![img.png](cmd/config/img.png)