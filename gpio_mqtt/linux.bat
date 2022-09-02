rem set GOARCH=arm
set GOOS=linux
rem set GOOS=windows


rem go tool dist install -v pkg/runtime
rem go install -v -a std

go build


rem set GOARCH=ARM64
rem go build -o aarch64_gpio_mqtt

set GOARCH=amd64
go build -o amd64_gpio_mqtt
