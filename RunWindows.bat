@echo off
:UniBuild
go build Uni.go
if %errorlevel% EQU 0 goto UniRun
goto End
:UniRun
"Uni.exe" -config ../UniConfig.inf
if %errorlevel% EQU 1 goto UniBuild
rem ExitCode 1 means to restart UniBot
:End
