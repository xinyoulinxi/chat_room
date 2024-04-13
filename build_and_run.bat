@echo off
cls
echo Building and running Go program...
go build -o out/web_server.exe
if %ERRORLEVEL% == 0 (
    echo Running program...
    .\out\web_server.exe
) else (
    echo Build failed.
)
