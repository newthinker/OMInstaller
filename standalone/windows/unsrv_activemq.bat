@echo uninstall activemq service .... Press Control-C to abort
@pause
@echo cd dir
cd $ONEMAP_HOME\services\activemq5.4.1\bin\win32
@echo exec uninstall
call $ONEMAP_HOME\services\activemq5.4.1\bin\win32\UninstallService.bat
@echo .