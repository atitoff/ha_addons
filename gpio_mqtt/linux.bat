set GOARCH=amd64
rem set GOARCH=arm
set GOOS=linux
rem set GOOS=windows

rem go tool dist install -v pkg/runtime
rem go install -v -a std

go build
rem go build -o gpio_mqtt_arm
