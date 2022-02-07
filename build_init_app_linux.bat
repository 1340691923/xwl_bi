 set goos=linux&& go build  -ldflags="-w -s" -o bin/linux/init_app cmd/init_app/main.go
echo "build success"
