@echo off
rem Builds RapTech Portable: tests, the launcher (put inside the manager), then the manager.
setlocal
cd /d "%~dp0"
go test ./... || exit /b 1
go build -trimpath -ldflags "-H windowsgui -s -w" -o assets\launcher.exe .\cmd\launcher || exit /b 1
copy /y assets\launcher.exe build\bin\RapPortable-Launcher.exe >nul 2>&1
wails build -clean -trimpath || exit /b 1
copy /y assets\launcher.exe build\bin\RapPortable-Launcher.exe >nul
echo.
echo Built build\bin\RapTech-Portable.exe
