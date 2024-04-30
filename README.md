# TalkRater_Bot
The [competitive task](https://www.notion.so/Telegram-835fdb2333be43efae71edd41362792f) 
of the [GolangConf 2024](https://cfp.golangconf.ru/) 
conference is a Telegram bot that simplifies the process of evaluating reports.

## Install project
1) install [docker engine](https://docs.docker.com/engine/install/), [docker compose](https://docs.docker.com/compose/install/), [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
2) Run bash script
```shell
git clone https://github.com/BaizhumanovAlisher/TalkRater_Bot.git
cd TalkRater_Bot

mkdir secrets
echo -n "<YOUR_API_TOKEN_FOR_USER>" | cat > ./secrets/tg_api_token_user.txt
echo -n "<YOUR_API_TOKEN_FOR_ADMIN>" | cat > ./secrets/tg_api_token_admin.txt
echo -n "<YOUR_DB_PASSWORD>" | cat > ./secrets/db_password.txt

export DB_PASSWORD_FILE="$(pwd)/secrets/db_password.txt"
export TG_API_TOKEN_USER_FILE="$(pwd)/secrets/tg_api_token_user.txt"
export TG_API_TOKEN_ADMIN_FILE="$(pwd)/secrets/tg_api_token_admin.txt"
```

## Run project
```shell
docker-compose up
```

## Run only talk_rater_bot container
```shell
docker-compose up talk_rater_bot
```

## About using GORM
This project contain GORM library, but there exist some critics in gopher community.

In context of this project, it is proved to use GORM. Reasons:

1) non a high-load
2) basic CRUD operations
3) speed up development
4) auto migration

[More information about using ORM in Go](https://youtu.be/MBfjQBDZqt8?si=I80cyqQxswjJCNg1)

## Telegram API Library
The `Telebot` library was selected because:
1) many star in [GitHub](https://github.com/tucnak/telebot)
2) documentation in [go.dev](https://pkg.go.dev/gopkg.in/telebot.v3)
3) MIT License

## Default time pattern from csv schedule
`21/07/2024 10:00:00` - MSK time Zone

## TODO List
- add monitoring with Prometheus
- add backup before removing all data from DB
- add check sql injection