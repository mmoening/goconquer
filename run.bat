@echo off
set GOPATH=%~dp0
set ERRORLEVEL=0

cd ./bin/account
start server.exe

cd ../game
start server.exe

cd ../web
start server.exe

cd ../../