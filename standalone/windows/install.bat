@echo off

set inputs=%~1
set log=%ONEMAP_HOME%\inst_standalone.log

echo "Begin installing OneMap" >>%log%
echo %inputs%	>>%log%
for %%p in (%inputs%) do (
    rem 注册环境变量
    if exist %ONEMAP_HOME%\bin\command\env.bat (
        call %ONEMAP_HOME%\bin\command\env.bat
        move /Y %ONEMAP_HOME%\bin\command\env.bat  %ONEMAP_HOME%\bin\command\enviroment.bat
        if %errorlevel% neq 0 (
            echo "Rename enviroment script failed" >>%log%
        )		

        echo "Regist successfully"	>>%log%
    )
        
    if /i not "%%~p"=="" (	
        rem 运行系统配置工具
        if /i "%%~p"=="sysconfig" (
            java -jar %ONEMAP_HOME%\config\SystemConfig\SystemConfig.jar %ONEMAP_HOME%\config\SystemConfig
            
            if %errorlevel% neq 0 (
                echo "System config failed"	>>%log%
                exit /b 2
            ) 
            echo "System config successfully"	>>%log%
        )
    
        setlocal EnableDelayedExpansion
        echo "Install server type:%%p"	
        if /i "%%~p"=="db" (
            set "ori=$ORACLE_DATA"
            set "new=%ORACLE_BASE%\oradata\%ORACLE_SID%"
            rem 执行sql脚本，首先获取ORACLE_BASE环境变量
            if exist "%ORACLE_BASE%\oradata\%ORACLE_SID%" (
                rem 首先对文件中的变量进行替换
                if exist %ONEMAP_HOME%\db\GeoShareManager\geoshare_platform.sql (
                    call :REPLACE %ONEMAP_HOME%\db\GeoShareManager\geoshare_platform.sql
                    if %errorlevel% neq 0 (
                        echo "Update manager script failed"		>>%log%
                        exit /b 6
                    )
                    %ORACLE_HOME%\bin\sqlplus.exe %ORACLE_SYSTEM_ACCOUNT%/%ORACLE_SYSTEM_PWD%@%ORACLE_SID% as sysdba @%ONEMAP_HOME%\db\GeoShareManager\geoshare_platform.sql <NUL
                    if %errorlevel% neq 0 (
                        echo "Create manager tablespace failed"		>>%log%
                        exit /b 7
                    )
                    %ORACLE_HOME%\bin\sqlplus.exe GEOSHARE_PLATFORM/admin@%ORACLE_SID% @%ONEMAP_HOME%\db\GeoShareManager\Manager_Table_Script.sql	<NUL
                    if %errorlevel% neq 0 (
                        echo "Create manager tables failed"		>>%log%
                        exit /b 8
                    )
                    %ORACLE_HOME%\bin\sqlplus.exe GEOSHARE_PLATFORM/admin@%ORACLE_SID% @%ONEMAP_HOME%\db\GeoShareManager\Manager_Table_Data.sql 	<NUL			
                    if %errorlevel% neq 0 (
                        echo "Import manager data failed"	>>%log%
                        exit /b 9
                    )
                    echo "Install manager database successfully"	>>%log%
                )
                if exist %ONEMAP_HOME%\db\Portal\geoshare_portal.sql (
                    call :REPLACE %ONEMAP_HOME%\db\Portal\geoshare_portal.sql
                    if %errorlevel% neq 0 (
                        echo "Update portal script failed"	>>%log%
                        exit /b 11
                    )
                    %ORACLE_HOME%\bin\sqlplus.exe %ORACLE_SYSTEM_ACCOUNT%/%ORACLE_SYSTEM_PWD%@%ORACLE_SID% as sysdba @%ONEMAP_HOME%\db\Portal\geoshare_portal.sql 	<NUL
                    if %errorlevel% neq 0 (
                        echo "Create portal tablespace failed"	>>%log%
                        exit /b 12
                    )
                    %ORACLE_HOME%\bin\sqlplus.exe GEOSHARE_PORTAL/admin@%ORACLE_SID% @%ONEMAP_HOME%\db\Portal\Portal_Table_Script.sql  <NUL
                    if %errorlevel% neq 0 (
                        echo "Create portal tables failed"	>>%log%
                        exit /b 13
                    )
                    %ORACLE_HOME%\bin\sqlplus.exe GEOSHARE_PORTAL/admin@%ORACLE_SID% @%ONEMAP_HOME%\db\Portal\Portal_Table_Data.sql  <NUL
                    if %errorlevel% neq 0 (
                        echo "Import portal data failed"	>>%log%
                        exit /b 14
                    )
                    echo "Install portal database successfully"	>>%log%
                )
                if exist %ONEMAP_HOME%\db\GeoCoding\geo_coding.sql (
                    call :REPLACE %ONEMAP_HOME%\db\GeoCoding\geo_coding.sql
                    if %errorlevel% neq 0 (
                        echo "Update geocoding script failed"	>>%log%
                        exit /b 16
                    )
                    %ORACLE_HOME%\bin\sqlplus.exe %ORACLE_SYSTEM_ACCOUNT%/%ORACLE_SYSTEM_PWD%@%ORACLE_SID% as sysdba @%ONEMAP_HOME%\db\GeoCoding\geo_coding.sql 	<NUL	
                    if %errorlevel% neq 0 (
                        echo "Create geocoding tablespace failed"	>>%log%
                        exit /b 17
                    )
                    %ORACLE_HOME%\bin\sqlplus.exe GEO_CODING/admin@%ORACLE_SID% @%ONEMAP_HOME%\db\GeoCoding\geo_coding_table.sql 		<NUL
                    if %errorlevel% neq 0 (
                        echo "Create portal tables failed"	>>%log%
                        exit /b 18
                    )
                    echo "Install geocoding database successfully"	>>%log%
                )
                if exist %ONEMAP_HOME%\db\GeoPortal\geo_portal.sql (
                    call :REPLACE %ONEMAP_HOME%\db\GeoPortal\geo_portal.sql
                    if %errorlevel% neq 0 (
                        echo "Update geoportal script failed"	>>%log%
                        exit /b 21
                    )
                    %ORACLE_HOME%\bin\sqlplus.exe %ORACLE_SYSTEM_ACCOUNT%/%ORACLE_SYSTEM_PWD%@%ORACLE_SID% as sysdba @%ONEMAP_HOME%\db\GeoPortal\geo_portal.sql 	<NUL
                    if %errorlevel% neq 0 (
                        echo "Create geoportal tablespace failed"	>>%log%
                        exit /b 22
                    )
                    %ORACLE_HOME%\bin\sqlplus.exe GEO_PORTAL/admin@%ORACLE_SID% @%ONEMAP_HOME%\db\GeoPortal\geo_portal_table.sql  	<NUL
                    if %errorlevel% neq 0 (
                        echo "Create geoportal tables failed"	>>%log%
                        exit /b 23
                    )
                    %ORACLE_HOME%\bin\sqlplus.exe GEO_PORTAL/admin@%ORACLE_SID% @%ONEMAP_HOME%\db\GeoPortal\geo_portal_data.sql  	<NUL
                    if %errorlevel% neq 0 (
                        echo "Import geoportal data failed"	>>%log%
                        exit /b 24
                    )
                    echo "Install geoportal database successfully"	>>%log%
                )	
                if exist %ONEMAP_HOME%\db\SubPlatform\geoshare_sub_platform.sql (
                    call :REPLACE %ONEMAP_HOME%\db\SubPlatform\geoshare_sub_platform.sql
                    if %errorlevel% neq 0 (
                        echo "Update sub script failed"	>>%log%
                        exit /b 26
                    )
                    %ORACLE_HOME%\bin\sqlplus.exe %ORACLE_SYSTEM_ACCOUNT%/%ORACLE_SYSTEM_PWD%@%ORACLE_SID% as sysdba @%ONEMAP_HOME%\db\SubPlatform\geoshare_sub_platform.sql 	<NUL
                    if %errorlevel% neq 0 (
                        echo "Create sub tablespace failed"	>>%log%
                        exit /b 27
                    )
                    %ORACLE_HOME%\bin\sqlplus.exe GEOSHARE_PLATFORM/admin@%ORACLE_SID% @%ONEMAP_HOME%\db\SubPlatform\Sub_Table_Script.sql  	<NUL
                    if %errorlevel% neq 0 (
                        echo "Create sub tables failed"	>>%log%
                        exit /b 28
                    )
                )
            )		
        )
        if /i "%%~p"=="gis" (
            rem 拷贝ojdbc库	
            if exist %AGS_HOME%\java\manager\config\security\lib (
                copy %ONEMAP_HOME%\db\Driver\ojdbc5.jar %AGS_HOME%\java\manager\config\security\lib
                copy %ONEMAP_HOME%\db\Driver\ojdbc5_g.jar  %AGS_HOME%\java\manager\config\security\lib
                
                if %errorlevel% neq 0 (
                    echo "Copy ojdbc package failed"	>>%log%
                    exit /b 31
                )
            )
        )
        if /i "%%~p"=="main" (
            rem nothing
        )
        if /i "%%~p"=="web" (
            rem nothing
        )
        if /i "%%~p"=="token" (
            rem nothing
        )
        if /i "%%~p"=="agent" (
            rem nothing
        )
        if /i "%%~p"=="msg" (
            rem nothing
        )
        echo "Install %%~p successfully"	>>%log%
    )
)
rem 注销让环境变量生效
shutdown.exe -l
goto :EOF

REM 对指定文件进行字符串替换
REM $0 {filename}	(input filename is quoted)
:REPLACE
set file=%1
rem echo %file%
if exist %file% (
    rem echo "%ori%"
    rem echo "%new%"
    
    rem delete the temp file
    if exist %file%_tmp.txt (
        del %file%_tmp.txt /f /s /q
    )
    
    for /f "delims=" %%i in ('findstr .* %file%') do (
        set "str=%%i"
        rem echo "The string is:"%str%
        rem deal with null line
        if not "!str!"=="" set "str=!str:%ori%=%new%!"
        >>%file%_tmp.txt echo.!str!
    )
    rem replace the original file
    move /Y %file%_tmp.txt  %file%
)
goto :EOF
