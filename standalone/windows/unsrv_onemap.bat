@echo uninstall OneMap service.... Press Control-C to abort
@pause
@echo cd dir
cd %ONEMAP_HOME%\Tomcat6.0.29\bin
@echo exec uninstall
call %ONEMAP_HOME%\Tomcat6.0.29\bin\Uninstall.bat
@echo .