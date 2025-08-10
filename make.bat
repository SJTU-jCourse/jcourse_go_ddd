@echo off
setlocal enabledelayedexpansion

set LOCAL_MODULE=jcourse_go

if "%1"=="" (
    echo Usage: %0 [lint^|migrate]
    echo   lint    - Run code formatting and import organization
    echo   migrate - Run database migrations
    exit /b 1
)

if "%1"=="lint" (
    echo Running go fmt...
    go fmt ./...
    
    echo Running goimports...
    for /f "delims=" %%f in ('dir /s /b *.go') do (
        goimports -local !LOCAL_MODULE! -w "%%f"
    )
    
    echo Running go mod tidy...
    go mod tidy
    
    echo Lint completed successfully!
    exit /b 0
)

if "%1"=="migrate" (
    echo Running database migrations...
    go run cmd/migrate/main.go
    exit /b 0
)

echo Unknown target: %1
echo Usage: %0 [lint^|migrate]
echo   lint    - Run code formatting and import organization
echo   migrate - Run database migrations
exit /b 1