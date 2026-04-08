@echo off
cd /d "%~dp0"

echo [1/3] Building Vue frontend...
pushd "%~dp0..\frontend"
call npm run build
if %ERRORLEVEL% neq 0 (
    echo [ERROR] npm run build failed.
    popd
    pause
    exit /b 1
)
popd

echo [2/3] Copying dist to backend\web...
if exist "%~dp0web" rmdir /S /Q "%~dp0web"
xcopy /E /Y /I "%~dp0..\frontend\dist" "%~dp0web\" >nul
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Failed to copy dist.
    pause
    exit /b 1
)

echo [3/3] Compiling Go for Linux/amd64 (prod embed)...
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -tags prod -ldflags="-s -w" -o "%~dp0silverlord" "%~dp0."
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Go build failed.
    pause
    exit /b 1
)

echo.
echo [OK] Output: silverlord  (Linux amd64, frontend embedded)
pause
