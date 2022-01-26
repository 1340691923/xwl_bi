set dir = %cd%

set args = "-configFileDir=%cd%\config"

gowatch  -p cmd/init_app/main.go -args=args -o bin/win/init_app.exe