@echo off

:: 用户级环境变量注册表项
::set regpath=HKEY_CURRENT_USER\Environment
:: 系统级环境变量注册表项
set regpath=HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment

echo 更新环境变量 Path
SETLOCAL ENABLEDELAYEDEXPANSION 
set "fqElement=C:\OneMap\Java\jdk1.6.26\bin"
echo java_home is:%fqElement%
@REM 将PATH环境变量中的路径转换成用空格分开的列表
set fpath="%PATH:;=" "%"
@REM 遍历PATH列表，删除指定的元素
for %%p in (%fpath%) do (
    @REM is not null
    if /i not "%%~p"=="" (
		echo %%p
		echo %%~p
    	@REM is this element NOT the one we want to remove?
    	if /i NOT "%%~p"=="%fqElement%" (
        	if _!tpath!==_ (set tpath=%%~p) else (set tpath=!tpath!;%%~p)
    	)
    )
)
echo %tpath%
set path=%tpath%
@REM 更新PATH环境变量
reg add "%regpath%" /v "PATH" /t REG_EXPAND_SZ /d "%PATH%;" /f
echo %PATH%

echo 删除环境变量 JAVA_HOME
reg delete "%regpath%" /v JAVA_HOME /f
echo.

echo 删除环境变量 CATALINA_HOME
reg delete "%regpath%" /v CATALINA_HOME /f
echo.

echo 删除环境变量 CATALINA_BASE
reg delete "%regpath%" /v CATALINA_BASE /f
echo.
 
echo 删除环境变量 CLASSPATH
reg delete "%regpath%" /v CLASSPATH /f
echo.

::echo 删除环境变量 JAVA_OPTS
::reg delete "%regpath%" /v JAVA_OPTS  /f
::echo.

echo 删除环境变量 ONEMAP_HOME
reg delete "%regpath%" /v ONEMAP_HOME /f
echo.

