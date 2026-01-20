@echo off
setlocal

set ROOT_DIR=%~dp0\..
set OUT_DIR=%ROOT_DIR%\dist\release\plugins
if not defined GO_BIN set GO_BIN=go
if not defined GOTOOLCHAIN set GOTOOLCHAIN=go1.25.6
set GOTOOLCHAIN=%GOTOOLCHAIN%

if not exist "%OUT_DIR%" mkdir "%OUT_DIR%"

%GO_BIN% build -buildmode=plugin -o "%OUT_DIR%\nowpayments.dll" "%ROOT_DIR%\plugins\nowpayments"
%GO_BIN% build -buildmode=plugin -o "%OUT_DIR%\coinbasecommerce.dll" "%ROOT_DIR%\plugins\coinbasecommerce"
%GO_BIN% build -buildmode=plugin -o "%OUT_DIR%\bitpay.dll" "%ROOT_DIR%\plugins\bitpay"
%GO_BIN% build -buildmode=plugin -o "%OUT_DIR%\coingate.dll" "%ROOT_DIR%\plugins\coingate"
%GO_BIN% build -buildmode=plugin -o "%OUT_DIR%\coinpayments.dll" "%ROOT_DIR%\plugins\coinpayments"

echo Plugins built in %OUT_DIR%
endlocal
