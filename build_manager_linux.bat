cmd /k "cd vue && npm run build:stage && cd .. && set goos=linux&&go build -o bin/linux/manager cmd/manager/main.go "
echo "build success"
