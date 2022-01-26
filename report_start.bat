set dir = %cd%

set args = "-configFileDir=%cd%\config"

gowatch  -p cmd/report_server/main.go -args=args -o bin/win/report_server.exe