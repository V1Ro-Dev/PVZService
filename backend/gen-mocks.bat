@echo off
setlocal

:: Путь к mockgen через go run
set MOCKGEN=go run github.com/golang/mock/mockgen

:: Пути
set DELIVERY_PATH=internal/delivery
set USECASE_PATH=internal/usecase
set REPOSITORY_PATH=internal/repository

echo Running tests with coverage...
for /f %%i in ('go list ./... ^| findstr /V "mocks"') do (
    go test %%i -coverprofile=coverage.out
)


echo.
echo Showing summarized coverage:
go tool cover -func=coverage.out

endlocal
pause
