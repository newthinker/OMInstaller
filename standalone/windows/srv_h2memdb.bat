@echo off
setlocal

@echo install H2MemDB service.... Press Control-C to abort
@echo cd dir
cd %ONEMAP_HOME%\services\H2CommonMemDB
@echo exec install
call %ONEMAP_HOME%\services\H2CommonMemDB\installH2MemDB.bat
net start H2CommonMemDB
@echo .

exit


