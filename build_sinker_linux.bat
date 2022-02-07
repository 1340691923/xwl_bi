set goos=linux&& go build  -ldflags="-w -s" -o bin/linux/sinker cmd/sinker/main.go
echo "build success"
