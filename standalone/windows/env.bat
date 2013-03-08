@echo off
:: �û�����������ע�����
::set regpath=HKEY_CURRENT_USER\Environment
:: ϵͳ����������ע�����
set regpath=HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment

rem set ONEMAP_HOME="C:\OneMap"

set javahome=%ONEMAP_HOME%\Java\jdk1.6.26
set activemqhome=%ONEMAP_HOME%\services\activemq5.4.1
set tomcathome=%ONEMAP_HOME%\Tomcat6.0.29
 
echo �½��������� JAVA_HOME=%javahome%
reg add "%regpath%" /v "JAVA_HOME" /d %javahome% /f
echo.

echo �½��������� CATALINA_HOME=%tomcathome%
reg add "%regpath%" /v "CATALINA_HOME" /d %tomcathome% /f
echo.

echo �½��������� CATALINA_BASE=%tomcathome%
reg add "%regpath%" /v "CATALINA_BASE" /d %tomcathome% /f
echo.
 
echo �½��������� CLASSPATH=.;%activemqhome%\lib;%javahome%\lib\tools.jar;%javahome%\lib\dt.jar
reg add "%regpath%" /v "CLASSPATH" /d .;%activemqhome%\lib;%javahome%\lib\tools.jar;%javahome%\lib\dt.jar /t REG_EXPAND_SZ /f
echo.

::echo �½��������� JAVA_OPTS=-Xms2048m -Xmx2048m -Duser.timezone=GMT+08 -XX:PermSize=256M -XX:MaxPermSize=1024M -Dsun.lang.ClassLoader.allowArraySyntax=true
::reg add "%regpath%" /v "JAVA_OPTS"  /d "-Xms2048m -Xmx2048m -Duser.timezone=GMT+08 -XX:PermSize=256M -XX:MaxPermSize=1024M -Dsun.lang.ClassLoader.allowArraySyntax=true" /f
::echo.

echo ׷�ӻ������� Path=%javahome%\bin
reg add "%regpath%" /v "Path" /t REG_EXPAND_SZ /d %javahome%\bin;"%Path%";  /f
echo.