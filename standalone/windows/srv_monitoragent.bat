@echo off
setlocal

@echo install MonitorAgent service.... Press Control-C to abort
@echo cd dir
cd %ONEMAP_HOME%\services\GeoShareMonitorAgent
@echo exec install
call %ONEMAP_HOME%\services\GeoShareMonitorAgent\InstallMonitorAgent.bat
net start monitoragent
@echo .

exit


