@echo uninstall H2MemDB service.... Press Control-C to abort
@pause
@echo cd dir
cd %ONEMAP_HOME%\services\H2CommonMemDB
@echo exec uninstall
net stop H2CommonMemDB
call %ONEMAP_HOME%\services\H2CommonMemDB\uninstallH2MemDB.bat
@echo .