# Telegram bot in go
## With only 2 libs ([godotenv](github.com/joho/godotenv) and [pq](github.com/lib/pq))
Commands
- /rnd        - random link
- /list       - all links
- https://... - save links
## How to start bot
1. Create .env
   ```
   TELEGRAM_TOKEN="your token from BotFather"
   DB_URL="postgres://..."
   ```
2. Build main.go (``` go build main.go ```)
3. Start main.exe
## Another storage
You can use file storage. You will not need DB_URL in .env, but this storage will require space in your disk.   
To use file storage:
1. Create folder in root of project
2. In main.go change *const storagePath* to name of your folder
3. In main.go add import "github.com/zumosik/telegram-go/storage/files" 
4. Comment this lines:
   ```
   storage, err := postgres.New(cfg.DatabaseURL)
   if err != nil {
       log.Fatalf("Error creating storage: %s", err)
	 }
   eventsProcessor := telegram.New(tgClient.New(tgBotHost, cfg.Token), storage)
   ```
5. Uncomment this line:
   ```
   eventsProcessor := telegram.New(tgClient.New(tgBotHost, cfg.Token), files.New(storagePath))
   ```
6. You will get code like this:
   ```
   eventsProcessor := telegram.New(tgClient.New(tgBotHost, cfg.Token), files.New(storagePath))
	// storage, err := postgres.New(cfg.DatabaseURL)
	// if err != nil {
	// 	log.Fatalf("Error creating storage: %s", err)
	// }
	// eventsProcessor := telegram.New(tgClient.New(tgBotHost, cfg.Token), storage)
   ```

![logo](https://github.com/zumosik/telegram-go/assets/86283476/a7ae023e-3048-4d62-b70e-2c9a6973065f)
