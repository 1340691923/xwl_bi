set dir = %cd%

set args = "-configFileDir=%cd%\config"

gowatch  -p cmd/manager/main.go -args=args  -o bin/win/manager.exe