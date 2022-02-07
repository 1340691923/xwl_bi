 set goos=linux&& go build  -ldflags="-w -s"  -o bin/linux/report_server cmd/report_server/main.go
echo "build success"
