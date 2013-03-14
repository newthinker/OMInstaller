@echo off
setlocal

@echo install OneMap service.... Press Control-C to abort
@echo cd dir
cd %ONEMAP_HOME%\Tomcat6.0.29\bin
@echo exec install
call %ONEMAP_HOME%\Tomcat6.0.29\bin\service.bat install
net start Tomcat6
@echo .

exit


