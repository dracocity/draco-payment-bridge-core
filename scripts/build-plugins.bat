@echo off
setlocal

set ROOT_DIR=%~dp0\..
set OUT_DIR=%ROOT_DIR%\plugins\dist

if not exist "%OUT_DIR%" mkdir "%OUT_DIR%"

go build -buildmode=plugin -o "%OUT_DIR%\nowpayments.dll" "%ROOT_DIR%\plugins\src\nowpayments"
go build -buildmode=plugin -o "%OUT_DIR%\coinbasecommerce.dll" "%ROOT_DIR%\plugins\src\coinbasecommerce"
go build -buildmode=plugin -o "%OUT_DIR%\bitpay.dll" "%ROOT_DIR%\plugins\src\bitpay"
go build -buildmode=plugin -o "%OUT_DIR%\coingate.dll" "%ROOT_DIR%\plugins\src\coingate"
go build -buildmode=plugin -o "%OUT_DIR%\coinpayments.dll" "%ROOT_DIR%\plugins\src\coinpayments"

echo Plugins built in %OUT_DIR%
endlocal
