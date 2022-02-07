cd vue && npm run build:stage && cd .. && set goos=linux&&go build  -ldflags="-w -s" -o bin/linux/manager cmd/manager/main.go
echo "build success"
