@echo off

:: �û�����������ע�����
::set regpath=HKEY_CURRENT_USER\Environment
:: ϵͳ����������ע�����
set regpath=HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment

echo ���»������� Path
SETLOCAL ENABLEDELAYEDEXPANSION 
set "fqElement=C:\OneMap\Java\jdk1.6.26\bin"
echo java_home is:%fqElement%
@REM ��PATH���������е�·��ת�����ÿո�ֿ����б�
set fpath="%PATH:;=" "%"
@REM ����PATH�б�ɾ��ָ����Ԫ��
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
@REM ����PATH��������
reg add "%regpath%" /v "PATH" /t REG_EXPAND_SZ /d "%PATH%;" /f
echo %PATH%

echo ɾ���������� JAVA_HOME
reg delete "%regpath%" /v JAVA_HOME /f
echo.

echo ɾ���������� CATALINA_HOME
reg delete "%regpath%" /v CATALINA_HOME /f
echo.

echo ɾ���������� CATALINA_BASE
reg delete "%regpath%" /v CATALINA_BASE /f
echo.
 
echo ɾ���������� CLASSPATH
reg delete "%regpath%" /v CLASSPATH /f
echo.

::echo ɾ���������� JAVA_OPTS
::reg delete "%regpath%" /v JAVA_OPTS  /f
::echo.

echo ɾ���������� ONEMAP_HOME
reg delete "%regpath%" /v ONEMAP_HOME /f
echo.

