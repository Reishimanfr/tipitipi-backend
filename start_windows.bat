@echo off
set DIST_PATH=.\src\backend\dist
set URL=http://localhost:2333
set MAX_ATTEMPTS=10
set WAIT_TIME=1
set attempt=1

@REM Create dist folder if it doesn't exist
if not exist "%DIST_PATH%" mkdir "%DIST_PATH%"

@REM Build server binary if not built
if not exist "%DIST_PATH%\server.exe" (
    @REM Cd into backend folder to build code
    cd ./src/backend

    @REM Build server binary
    go build -o ".\dist\server.exe\"
    cd ../../
)

@REM Start server
start /B "%DIST_PATH%\server.exe"

@REM Wait for server to respond or exit
:LOOP
if %attempt% leq %MAX_ATTEMPTS% (
    echo Attempt %attempt%/%MAX_ATTEMPTS%...

    powershell -Command "try { $response = Invoke-WebRequest -Uri '%URL%' -Method Head -TimeoutSec 5; if ($response.StatusCode -eq 200) { exit 0 } } catch { exit 1 }"
    if %errorlevel% equ 0 (
        echo Server is up!
        @REM Start up vite if server responded
        bunx vite
        exit /b 0
    ) else (
        echo Attempt %attempt% failed. Retrying in %WAIT_TIME% seconds...
        timeout /t %WAIT_TIME% /nobreak >nul
    )

    set /a attempt=attempt+1
    goto LOOP
)

echo Failed to reach the server after %MAX_ATTEMPTS% attempts.
exit /b 1
