# TalkRater_Bot
The [competitive task](https://www.notion.so/Telegram-835fdb2333be43efae71edd41362792f) 
of the [GolangConf 2024](https://cfp.golangconf.ru/) 
conference is a Telegram bot that simplifies the process of evaluating reports.

## About using GORM
This project contain GORM library, but there exist some critics in gopher community.

In context of this project, it is proved to use GORM. Reasons:

1) non a high-load
2) basic CRUD operations
3) speed up development
4) auto migration

[More information about using ORM in Go](https://youtu.be/MBfjQBDZqt8?si=I80cyqQxswjJCNg1)

## Telegram API Library
The `telegram-bot-api` library was selected because:
1) ready-made library
2) the largest number of stars in [GitHub](https://github.com/go-telegram-bot-api/telegram-bot-api)
3) documentation in [go.dev](https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5)
