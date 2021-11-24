# Инструкции по использованию курсового проекта

------

## Начало работы
Для начала нужно подготовить несколько ENV переменных:
* TOKEN_PATH_PRIVATE -- записывается путь к файлу, где лежит private api токен
* TOKEN_PATH_PUBLIC -- записывается путь к файлу, где лежит public api токен
* WS_URL -- URL веб сокета, откуда будут тянуться свечки
* TG_BOT_TOKEN -- путь к файлу, где лежит токен телеграм бота

Для запуска курсового проекта, нужно запустить файл `main.go` (находится в корне папки проекта)

## Взаимодействие
Все взаимодействие с приложением происходит через телеграм [бота](http://t.me/Tfs_CP_bot)