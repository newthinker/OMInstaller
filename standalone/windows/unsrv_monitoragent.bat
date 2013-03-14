@echo uninstall MonitorAgent service.... Press Control-C to abort
@pause
@echo cd dir
cd %ONEMAP_HOME%\services\GeoShareMonitorAgent
@echo exec uninstall
net stop monitoragent
call %ONEMAP_HOME%\services\GeoShareMonitorAgent\UninstallMonitorAgent.bat
@echo .