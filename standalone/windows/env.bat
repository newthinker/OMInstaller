@echo off
:: 用户级环境变量注册表项
::set regpath=HKEY_CURRENT_USER\Environment
:: 系统级环境变量注册表项
set regpath=HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment

rem set ONEMAP_HOME="C:\OneMap"

set javahome=%ONEMAP_HOME%\Java\jdk1.6.26
set activemqhome=%ONEMAP_HOME%\services\activemq5.4.1
set tomcathome=%ONEMAP_HOME%\Tomcat6.0.29
 
echo 新建环境变量 JAVA_HOME=%javahome%
reg add "%regpath%" /v "JAVA_HOME" /d %javahome% /f
echo.

echo 新建环境变量 CATALINA_HOME=%tomcathome%
reg add "%regpath%" /v "CATALINA_HOME" /d %tomcathome% /f
echo.

echo 新建环境变量 CATALINA_BASE=%tomcathome%
reg add "%regpath%" /v "CATALINA_BASE" /d %tomcathome% /f
echo.
 
echo 新建环境变量 CLASSPATH=.;%activemqhome%\lib;%javahome%\lib\tools.jar;%javahome%\lib\dt.jar
reg add "%regpath%" /v "CLASSPATH" /d .;%activemqhome%\lib;%javahome%\lib\tools.jar;%javahome%\lib\dt.jar /t REG_EXPAND_SZ /f
echo.

::echo 新建环境变量 JAVA_OPTS=-Xms2048m -Xmx2048m -Duser.timezone=GMT+08 -XX:PermSize=256M -XX:MaxPermSize=1024M -Dsun.lang.ClassLoader.allowArraySyntax=true
::reg add "%regpath%" /v "JAVA_OPTS"  /d "-Xms2048m -Xmx2048m -Duser.timezone=GMT+08 -XX:PermSize=256M -XX:MaxPermSize=1024M -Dsun.lang.ClassLoader.allowArraySyntax=true" /f
::echo.

echo 追加环境变量 Path=%javahome%\bin
reg add "%regpath%" /v "Path" /t REG_EXPAND_SZ /d %javahome%\bin;"%Path%";  /f
echo.