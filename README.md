# Телеграм бот

Телеграм бот написанный на Go, бот написанный исключительно с помощью стандартной библиотеки Go (исключение драйвер для sqlite, но так же есть поддержка хранения данных в файловой системе).
Функционал: данный бот принимает какие-либо ссылки, например ссылка на статью которую сейчас вы не хотите читать, но собираетесь в обозримом будущем и не хотите ее потерять, бот сохраняет у себя ссылку и специальной командой выдает вам рандомную ссылку

---

## Установка
```bash
git clone https://github.com/ladev74/saves-given-link-bot
cd saves-given-link-bot
go build -o tg-bot
./tg-bot -tg-bot-token 'your token for telegram'
```

```text
Токен для бота можно получить в самом телеграме у бота BotFather
```