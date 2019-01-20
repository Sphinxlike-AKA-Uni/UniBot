@echo off
C:\go\bin\go build -buildmode=exe Uni.go
if %errorlevel% EQU 0 (
"!RunBot.bat"
) else (
echo Errorlevel: %errorlevel%
Pause
)
