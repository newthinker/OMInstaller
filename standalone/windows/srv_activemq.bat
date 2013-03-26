@echo off
setlocal

@echo install activemq service .... Press Control-C to abort
@echo cd dir
cd %ONEMAP_HOME%\services\activemq5.4.1\bin\win32
@echo exec install
call %ONEMAP_HOME%\services\activemq5.4.1\bin\win32\InstallService.bat
net start ActiveMQ
@echo .

exit


