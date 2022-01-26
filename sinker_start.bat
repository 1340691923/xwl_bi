set dir = %cd%

set args = "-configFileDir=%cd%\config"

gowatch  -p cmd/sinker/main.go -args=args -o bin/win/sinker.exe