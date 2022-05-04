# IP-TELEGRAM-BOT

## INSTALATION

Prerequisites:
- docker
- git

Simply clone this repo:

```
git clone https://github.com/faked86/ip-telegram-bot.git
cd ip-telegram-bot
```

## USAGE

1. Create `.env` file in project directory.
2. Place your telegram bot token in it: `TOKEN=YOUR_TOKEN`.
3. Run `docker-compose up` in terminal.

### Telegram bot

#### For regular users:
- Send IP address to bot to get info about it.
- Command `/start` to register you in database.
- Command `/unique` to get info about your unique requests.

#### For admins:
- Command `/spam <message>` to send `<message>` to all users in our database.
- Command `/admin <user>` to make `<user>` admin or make regular user if `<user>` is already admin.
- Command `/history <user>` to check `<user>`'s request history.

#### How to make first admin:
1. Register in database by running `/start` in telegram chat.
2. Connect to database via tool like pgAdmin 4:
      - POSTGRES USER: pg
      - POSTGRES PASSWORD: pass
      - POSTGRES DB: crud
3. Run query `UPDATE public.users SET admin = true WHERE username = '<your username>'`


### API Server

- `GET    /get_users` - Get all users.
- `GET    /get_users?id=<telegram id>` - Get user by `<telegram id>`.
- `GET    /get_history_by_tg?id=<telegram id>` - Get user history by `<telegram id>`.
- `DELETE /delete_from_history/<request_id>` - Delete request from history by `<request_id>`.
 You can get `<request_id>` on `/get_history_by_tg?id=<telegram id>` page.