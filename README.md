# TalkRater_Bot
The [competitive task](https://www.notion.so/Telegram-835fdb2333be43efae71edd41362792f) 
of the [GolangConf 2024](https://cfp.golangconf.ru/) 
conference is a Telegram bot that simplifies the process of evaluating reports.

## Install and run project
1) install [docker engine](https://docs.docker.com/engine/install/), [docker compose](https://docs.docker.com/compose/install/), [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
2) `git clone https://github.com/BaizhumanovAlisher/TalkRater_Bot.git`
3) `cd TalkRater_Bot`
4) `docker-compose up`

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
