@echo off
set GOPATH=%~dp0
set ERRORLEVEL=0

echo Building project...
go build -i -o ./bin/account/server.exe -v ./src/account/main
go build -i -o ./bin/game/server.exe -v ./src/game/main
go build -i -o ./bin/web/server.exe -v ./src/web/main
echo Build completed.