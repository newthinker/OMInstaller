#!/bin/bash

####################################################################################################
# DISP: Update version of install bash script for OneMap Linux version
# AUTH: Michael.Cho
# TIME: 2012/10/17
# EMAIL: zuow11@gmail.com
####################################################################################################

####################################################################################################
# consts and virables
####################################################################################################
MA="monitoragent"							# MonitorAgent service in onemap
MQ="activemq"								# Activemq service in onemap
H2="h2memdb"								# H2memdb service in onemap
ONEMAP="onemap"								# Onemap service 

#CONTAINER_HOME={$ONEMAP_HOME}"/"{$CONTAINER_NAME}{$CONTAINER_VERSION}	# Web server container home

####################################################################################################
# public functions
####################################################################################################
# Copy the public files of all module
#
# @RET : No return code
function copyPublic(){
	if ( ! [ -d $ONEMAP_HOME ] );then
		#mkdir -p $ONEMAP_HOME
        return 1
	fi

	# copy files
#	cp -r OneMap/java $ONEMAP_HOME
#	cp -r OneMap/config $ONEMAP_HOME
#	cp -r OneMap/temp $ONEMAP_HOME
#	cp -r OneMap/logs $ONEMAP_HOME
}

# @Discription: Regist the system enviroment
# 
function regEnv(){
	# OneMap
	if ( ! [ -d ${ONEMAP_HOME} ] );then 
		echo "ONEMAP_HOME directory isn't exist, please check!"
		return 1
	fi	
	echo "export ONEMAP_HOME=$ONEMAP_HOME" >> /etc/profile

	# JDK
	if ( ! [ -d ${JAVA_HOME} ] );then 
		echo "ONEMAP_HOME directory isn't exist, please check!"
		return 2
	fi	
	echo "JAVA_HOME=$ONEMAP_HOME/java/jdk1.6.0_26" >> /etc/profile
	echo "JRE_HOME=\$JAVA_HOME/jre" >> /etc/profile
	echo "CLASSPATH=\$CLASSPATH:.:\$JAVA_HOME/lib/dt.jar:\$JAVA_HOME/lib/tools.jar"	>> /etc/profile
	echo "PATH=\$PATH:\$JAVA_HOME/bin:\$JRE_HOME/bin"	>> /etc/profile
	echo "export JAVA_HOME JRE_HOME CLASSPATH PATH" >> /etc/profile
	
	# register the web server container
#	if ( ! [ -d ${CONTAINER_HOME} ] );then 
#		echo "ONEMAP_HOME directory isn't exist, please check!"
#		return 3
#	fi		
	local temp1=$(echo $CONTAINER_NAME | tr [a-z] [A-Z])
    local sub="TOMCAT"
    strstr $tmep1 $sub
	if test $? -eq 0
	then
		echo "export CATALINA_BASE=$CONTAINER_HOME" >> /etc/profile		
		echo "export CATALINA_HOME=$CONTAINER_HOME" >> /etc/profile		
    else
		echo "export MW_HOME=$CONTAINER_HOME" >> /etc/profile		
	fi
	
	# ulimit
	echo "ulimit -n $ULIMIT_NUM" >> /etc/profile
	
	source /etc/profile
	
	return 0
}

# @Discription: registration a system service
# @Arguments: regService SERVICE_NAME
# @RET: 0 => registration service successfully
#       !0 => registration failed
function regService(){
	if [ $# -eq 0 ];then
		echo "ERROR: Please input the registration service name"
		return 1
	elif [ $# -ne 1 ];then
		echo "ERROR: Please input only one service name per time"
		return 1
	fi

	srvName=$1
	cp ${ONEMAP_HOME}/bin/service/$srvName /etc/init.d/
	
	# registration the service
	chmod +x /etc/init.d/$srvName
	chkconfig --add $srvName
	RET=$?
	if ( ! [ $RET -eq 0 ] );then
		echo "ERROR: Registration service $srvName failed"
		return 2
	fi
	# set the status
	chkconfig $srvName on
	RET=$?
	if ( ! [ $RET -eq 0 ] );then
		echo "ERROR: Set the status of service $srvName failed"
		return 2
	fi
	# start the server
	#service $srvName start
#    /etc/init.d/${srvName} start &
#	if ( ! [ $RET -eq 0 ] );then
#		echo "ERROR: Start service $srvName failed"
#		return 2
#	fi
}

# Discription: Check a server running or not
# Arguments: chkServer SERVER_NAME
# @RET: 0 => the service is running
# 		1 => the service isn't running
#		2 => the service isn't installed
#		3 => others
function chkServer(){
	srvName=$1
	
	if test -f "/etc/init.d/$srvName"	
	then
		if test -f "/var/run/${srvName}.pid"
		then
			if test -z "$(cat /var/run/$srvName)"
			then 
				echo "INFO: The pid file of Service $srvName is null"
				return 3
			fi
			
			if test ps -p $(cat /var/run/$srvName) >/dev/null
			then
				echo "INFO: Server $srvName(pid: cat /var/run/${srvName}.pid)"
				return 0
			else
				echo "INFO: Server $srvName has terminated unexpectedly"
				return 1
			fi
		else
			echo "INFO: Server $srvName isn't running"
			return 1
		fi
	else
		echo "INFO: Server $srvName hasn't installed"
		return 2
	fi

}

# @Discription: Check the oracle database
# 
# @RET: 0 => good status
#		!0 => bad status, something wrong with the oracle installation
function checkOracle(){
	# check oracle account
	# local count=0
	# local flag=0
	# while [ $count -lt 3 ]
	# do
		# #check the oracle user and the oracle directory
		# read -p "Input the oracle account name(oracle): " oraUser
		# if test -z $oraUser
		# then
			# oraUser="oracle"
		# fi
		
		# RET="`grep \"^$oraUser\" /etc/passwd | wc -l`"
		# if test $RET -eq 1
		# then
			# flag=1
			# break
		# fi
		# echo "ERROR: Input error, please again"
	# done
	
	# if test $flag -eq 0
	# then
		# echo "ERROR: Input oracle account error and return"
		# return 1
	# fi
	
	if test -z ${ORCL_ACCOUNT}
	then
		echo "ERROR:Input oracle system account(${ORCL_ACCOUNT}) is null, please check!"
		return 1
	fi
	
	checkAccount ${ORCL_ACCOUNT}
	RET=$?
	if test $RET -ne 0
	then 
		echo "ERROR:Oracle system account(${ORCL_ACCOUNT}) doesn't existed, please check!"
		return 2
	fi
	
	# check oracle home
	export `su - "${ORCL_ACCOUNT}" -c "env | grep ^ORACLE_HOME "`
	if test -z $ORACLE_HOME
	then
		echo "ERROR:ORACLE_HOME enviroment($ORACLE_HOME) isn't exported, please check!"
		return 3
	elif ( ! [ -d $ORACLE_HOME ] ) 
	then
		echo "ERROR:Oracle home directory($ORACLE_HOME) does not existed, please check!"
		return 3
	fi

    # get the oracle base 
    ORACLE_BASE=`grep "<PROPERTY NAME=\"ORACLE_BASE\"" $ORACLE_HOME/inventory/ContentsXML/oraclehomeproperties.xml|awk '{print $3}'|awk -F\" '{print $2}'`
    if test -z $ORACLE_BASE
    then
        echo "ERROR: Get ORACLE_BASE failed and please check!"
        return 3
    fi
    export ORACLE_BASE
    ORACLE_DATA=${ORACLE_BASE}"/oradata/"${ORACLE_SID}
    export ORACLE_DATA

	# oracle sid
	export `su - "${ORCL_ACCOUNT}" -c "env | grep ^ORACLE_SID "`
	if [ -z $ORACLE_SID -o -z ${ORACLE_SID} ]
	then
		echo "ERROR:Oracle SID doesnot exist, please check!"
		return 4
	elif test $ORACLE_SID != ${ORACLE_SID}
	then
		echo "ERROR:System oracle sid($ORACLE_SID) isn't the same as the input oracle sid(${ORACLE_SID}), please check!"
		return 4
	fi	
}

# Discription: Check the ArcGIS Server
# 
# @RET: 0 => good status
#		!0 => bad status, something wrong with the AGS installation
function checkAGS() {
	# local count=0
	# local flag=0
	# while [ $count -lt 3 ]
	# do
		##check the AGS home directory
		# read -p "Input the ArcGIS Server home directory(/opt/arcgis/server10.0?): " AGS_HOME
		# if test -z $AGS_HOME
		# then
			# AGS_HOME="/opt/arcgis/server10.0"
		# fi
		
		# if test -d $AGS_HOME
		# then
			# flag=1
			# break
		# fi
		# echo "ERROR: Input ArcGIS Server home doesnot exist, please again"
	# done
	
	# Check AGS's home directory
	test -d ${AGS_HOME}
    if test $? -ne 0
	then
		echo "ERROR: Input ArcGIS Server home directory wrong and return"
		return 1
	fi
	
	return 0
}

# @Discription: Print the introduction information
# @RET: No return code
function printIntro(){

	echo ""
	echo "####################################################"
	echo ""
	echo "              OneMap installation tool              "
	echo ""
	echo "####################################################"
}

function printUsage(){
	echo "Usage: $0 {agent|db|gis|main|web|token|msg|all}"
	echo ""
	echo "Options:"
	echo "agent: install monitoragent server"
	echo "db: install database server"
	echo "gis: install gis server"
	echo "main: install maintenance server"
	echo "web: install web server"
	echo "token: install token server"
	echo "msg: install message server"
	echo "all: install all above servers"
	echo ""
	echo "Notice: When you input 'all' options, you needn't input others."
}

# @Discription:Check the account whether exist. 
# @Arguments: account_name 
# @RET: 0, success; !0, failed
# @LOG: [michael, 2012/10/17]
function checkAccount() {
	if test $# -lt 1
	then
		echo "Check account failed and please input the account name."
		echo "Usage: $0 {$account_name}"
		exit 1
	fi
	
	local account_name=$1
	# local account_pwd=$2

	# check the account whether exist
	local flag=`cat /etc/passwd | grep $account_name | wc -l`
	if test ${flag} -eq 1
	then
		echo "The $account_name has existed and continue."
	else	 
		echo "The $account_name has not existed!"
		return 1
		# first create $account_name group
		# groupadd $account_name
		# then create OneMap account
		# useradd -g $account_name -G root  -d /home/$account_name  $account_name
		# set password of One
		# echo $account_pwd | passwd --stdin $account_name
	fi
	
	return 0
}

####################################################################################################
# Install variables servers
####################################################################################################
# Discription: Install MonitorAgent Server
# @RET: 0 => Install successfully
#		!0 => Install failed
function instSrvMonitor(){
	echo "INFO: Start installing MonitorAgent Server"
	
	# if ( ! [ -d $ONEMAP_HOME/services ] );then
		# mkdir -p $ONEMAP_HOME/services
	# fi
	
	# copy files
##	copyPublic
##	cp -r OneMap/services/GeoShareMonitorAgent $ONEMAP_HOME/services/
	
	# set the enviroment
	regEnv
	RET=$?
	if [ $RET -ne 0 ];then
		echo "ERROR: Failed to register the enviroment parameters."
		return 1
	fi

    # exec SysConfig.sh
    ${ONEMAP_HOME}/config/SystemConfig/SysConfig.sh
    RET=$?
    if test $RET -ne 0 
    then
        echo "ERROR: Failed to config the OneMap"
        return 2
    fi

	# registration the service
	regService "$MA"
	RET=$?
	if [ $RET -ne 0 ];then
		echo "ERROR: Failed to install MonitorAgent Server"
		return 2
	fi
	
	echo "INFO: Install MonitorAgent Server successfully"
	
	return 0
}

# Discription: Install Database Server
# @RET: 0 => Install successfully
#		!0 => Install failed
function instSrvDB(){ 
	echo "INFO: Start installing Database Server"
	
	# first, check whether install MoniterAgent
	chkServer "$MA"
	RET=$?
	if test $RET -eq 2 
	then
		instSrvMonitor
	fi
	
	# check
	checkOracle
	RET=$?
	if test $RET -ne 0
	then
		echo "ERROR: Something wrong with oracle database, please check!"
		return 1
	fi
	
	# copy files
	#cp -r OneMap/db $ONEMAP_HOME/
	
	# update the ORACLE_HOME path of SQL scripts
	sed -i "s:\$ORACLE_DATA:$ORACLE_DATA:g"  ${ONEMAP_HOME}/db/GeoShareManager/geoshare_platform.sql
	sed -i "s:\$ORACLE_DATA:$ORACLE_DATA:g"  ${ONEMAP_HOME}/db/Portal/geoshare_portal.sql
	sed -i "s:\$ORACLE_DATA:$ORACLE_DATA:g"  ${ONEMAP_HOME}/db/GeoCoding/geo_coding.sql
	sed -i "s:\$ORACLE_DATA:$ORACLE_DATA:g"  ${ONEMAP_HOME}/db/GeoPortal/geo_portal.sql
	
	# create the database and import data
	su - ${ORCL_ACCOUNT} -c "/bin/sh -c \"echo exit | $ORACLE_HOME/bin/sqlplus ${ORACLE_SYSTEM_ACCOUNT}/${ORACLE_SYSTEM_PWD} as sysdba @$ONEMAP_HOME/db/GeoShareManager/geoshare_platform.sql \"" 
	su - ${ORCL_ACCOUNT} -c "/bin/sh -c \"echo exit | $ORACLE_HOME/bin/sqlplus ${MANAGER_USER}/${MANAGER_PWD} @$ONEMAP_HOME/db/GeoShareManager/Manager_Table_Script.sql \"" 
	su - ${ORCL_ACCOUNT} -c "/bin/sh -c \"echo exit | $ORACLE_HOME/bin/sqlplus ${MANAGER_USER}/${MANAGER_PWD} @$ONEMAP_HOME/db/GeoShareManager/Manager_Table_Data.sql \"" 
	
	su - ${ORCL_ACCOUNT} -c "/bin/sh -c \"echo exit | $ORACLE_HOME/bin/sqlplus ${ORACLE_SYSTEM_ACCOUNT}/${ORACLE_SYSTEM_PWD} as sysdba @$ONEMAP_HOME/db/Portal/geoshare_portal.sql \""
	su - ${ORCL_ACCOUNT} -c "/bin/sh -c \"echo exit | $ORACLE_HOME/bin/sqlplus ${PORTAL_USER}/${PORTAL_PWD} @$ONEMAP_HOME/db/Portal/Portal_Table_Script.sql \"" 
	su - ${ORCL_ACCOUNT} -c "/bin/sh -c \"echo exit | $ORACLE_HOME/bin/sqlplus ${PORTAL_USER}/${PORTAL_PWD} @$ONEMAP_HOME/db/Portal/Portal_Table_Data.sql \"" 

	su - ${ORCL_ACCOUNT} -c "/bin/sh -c \"echo exit | $ORACLE_HOME/bin/sqlplus ${ORACLE_SYSTEM_ACCOUNT}/${ORACLE_SYSTEM_PWD} as sysdba @$ONEMAP_HOME/db/GeoCoding/geo_coding.sql \""
	su - ${ORCL_ACCOUNT} -c "/bin/sh -c \"echo exit | $ORACLE_HOME/bin/sqlplus ${GEOCODING_USER}/${GEOCODING_PWD} @$ONEMAP_HOME/db/GeoCoding/geo_coding_table.sql \"" 
	
	su - ${ORCL_ACCOUNT} -c "/bin/sh -c \"echo exit | $ORACLE_HOME/bin/sqlplus ${ORACLE_SYSTEM_ACCOUNT}/${ORACLE_SYSTEM_PWD} as sysdba @$ONEMAP_HOME/db/GeoPortal/geo_portal.sql \""
	su - ${ORCL_ACCOUNT} -c "/bin/sh -c \"echo exit | $ORACLE_HOME/bin/sqlplus ${GEOPORTAL_USER}/${GEOPORTAL_PWD} @$ONEMAP_HOME/db/GeoPortal/geo_portal_table.sql \"" 	
	su - ${ORCL_ACCOUNT} -c "/bin/sh -c \"echo exit | $ORACLE_HOME/bin/sqlplus ${GEOPORTAL_USER}/${GEOPORTAL_PWD} @$ONEMAP_HOME/db/GeoPortal/geo_portal_data.sql \"" 	
	
	echo "INFO: Finish installing Database server"
			
	return 0
}

# Discription: Install GIS Server
# @RET: 0 => Install successfully
#		!0 => Install failed
function instSrvGIS(){
	echo "INFO: Begin installing GIS Server..."
	
	checkAGS
	RET=$?
	if test $RET -ne 0
	then
		echo "ERROR: Install GIS Server failed and return"
		return 1
	fi

	# first, check whether install MoniterAgent
	chkServer "$MA"
	RET=$?
	if test $RET -eq 2 
	then
		instSrvMonitor
	fi
	
	# copy jdbc driver files
	cp ${ONEMAP_HOME}/db/Driver/ojdbc5.jar ${ONEMAP_HOME}/db/Driver/ojdbc5_g.jar  ${AGS_HOME}/java/manager/config/security/lib/
	chmod 777 ${AGS_HOME}/java/manager/config/security/lib/ojdbc5*.jar
	
	if ( ! [ -d ${ONEMAP_HOME}/arcgis/license ] );then
		mkdir -p ${ONEMAP_HOME}/arcgis/license
	fi
		
	# check whether installed web container
	chkServer "$ONEMAP"
	RET=$?
	if test $RET -eq 0	# the service is running then terminal it first
	then
		service $ONEMAP stop
	elif test $RET -eq 2 
	then
		# cp -r OneMap/WLS1035 $ONEMAP_HOME
		
		# if ( ! [ -d ${ONEMAP_HOME}/WLS1035/onemap/onemap/webapps ] );then 
			# mkdir -p $ONEMAP_HOME/WLS1035/onemap/onemap/webapps
		# fi

		# cp -r OneMap/webapps/arcgis*  OneMap/webapps/LogServer  OneMap/webapps/Remote*  OneMap/webapps/ROOT $ONEMAP_HOME/WLS1035/onemap/onemap/webapps	
			
		regService "$ONEMAP"
		RET=$?
		if [ $RET -ne 0 ];then
			echo "ERROR: Failed to install onemap Server"
			return 1
		fi	
	# else
		# cp -r OneMap/webapps/arcgis  OneMap/webapps/LogServer  OneMap/webapps/Remote* $ONEMAP_HOME/WLS1035/onemap/onemap/webapps	
			
		# if ( ! [ -d $ONEMAP_HOME/WLS1035/onemap/onemap/webapps/ROOT ] );then
			# cp -r OneMap/webapps/ROOT	$ONEMAP_HOME/WLS1035/onemap/onemap/webapps
		# fi
		
		# if ( ! [ -d ${ONEMAP_HOME}/WLS1035/onemap/onemap/webapps/arcgis_js_api ] );then
			# cp -r OneMap/webapps/arcgis_js_api	${ONEMAP_HOME}/WLS1035/onemap/onemap/webapps
		# fi
	fi
	
	echo "INFO: Finish installing GIS Server"
	return 0
}

# Discription: Install Maintenance Server
# @RET: 0 => Install successfully
#		!0 => Install failed
function instSrvMain(){
	# first, check whether install MoniterAgent
	chkServer "$MA"
	RET=$?
	if test $RET -eq 2 
	then
		instSrvMonitor
	fi
	
	# second, check whether install h2memdb service
	chkServer "$H2"
	RET=$?
	if test $RET -eq 2 
	then
		# cp -r OneMap/services/H2CommonMemDB $ONEMAP_HOME/services/
		regService "$H2"
		RET=$?
		if test $RET -ne 0
		then
			echo "ERROR: Registration Service $H2 failed and please check"
		fi
	fi
	
	# copy files
	if ( ! [ -d $ONEMAP_HOME/arcgis/license ] );then
		mkdir -p $ONEMAP_HOME/arcgis/license
	fi	
	
	# check whether installed web container
	chkServer "$ONEMAP"
	RET=$?
	if test $RET -eq 0	# the service is running then terminal it first
	then
		service $ONEMAP stop
	elif test $RET -eq 2 
	then 
		# cp -r OneMap/WLS1035 $ONEMAP_HOME
		# if ( ! [ -d ${ONEMAP_HOME}/WLS1035/onemap/onemap/webapps ] );then 
			# mkdir -p $ONEMAP_HOME/WLS1035/onemap/onemap/webapps
		# fi

		# if ( ! [ -d $ONEMAP_HOME/WLS1035/onemap/onemap/webapps/ROOT ] );then
			# cp -r OneMap/webapps/ROOT	$ONEMAP_HOME/WLS1035/onemap/onemap/webapps
		# fi
		
		# cp -r OneMap/webapps/GeoShareManager  OneMap/webapps/BufferAnalyst	OneMap/webapps/ROOT ${ONEMAP_HOME}/WLS1035/onemap/onemap/webapps
		
		regService "$ONEMAP"
		RET=$?
		if [ $RET -ne 0 ];then
			echo "ERROR: Failed to install onemap Server"
			return 1
		fi	
	# else
		# cp -r OneMap/webapps/GeoShareManager  OneMap/webapps/BufferAnalyst 	${ONEMAP_HOME}/WLS1035/onemap/onemap/webapps

		# if ( ! [ -d $ONEMAP_HOME/WLS1035/onemap/onemap/webapps/ROOT ] );then
			# cp -r OneMap/webapps/ROOT	$ONEMAP_HOME/WLS1035/onemap/onemap/webapps
		# fi
	fi
	
	echo "INFO: Finish installing Maintenance Server"
	return 0
}

# Discription: Install WEB Server
# @RET: 0 => Install successfully
#		!0 => Install failed
function instSrvWEB(){
	# first, check whether install MoniterAgent
	chkServer "$MA"
	RET=$?
	if test $RET -eq 2 
	then
		instSrvMonitor
	fi
	
	# check whether installed web container
	chkServer "$ONEMAP"
	RET=$?
	if test $RET -eq 0	# the service is running then terminal it first
	then
		service $ONEMAP stop
	elif test $RET -eq 2 
	then 
		# cp -r OneMap/WLS1035 $ONEMAP_HOME
		# if ( ! [ -d ${ONEMAP_HOME}/WLS1035/onemap/onemap/webapps ] );then 
			# mkdir -p $ONEMAP_HOME/WLS1035/onemap/onemap/webapps
		# fi
		
		# cp -r OneMap/webapps/Portal  OneMap/webapps/RIA_APP_WIZARD OneMap/webapps/arcgis_js_api ${ONEMAP_HOME}/WLS1035/onemap/onemap/webapps
		
		regService "$ONEMAP"
		RET=$?
		if [ $RET -ne 0 ];then
			echo "ERROR: Failed to install onemap Server"
			return 1
		fi	
	# else
		# cp -r OneMap/webapps/Portal  OneMap/webapps/RIA_APP_WIZARD ${ONEMAP_HOME}/WLS1035/onemap/onemap/webapps
			
		# if ( ! [ -d ${ONEMAP_HOME}/WLS1035/onemap/onemap/webapps/arcgis_js_api ] );then
			# cp -r OneMap/webapps/arcgis_js_api	${ONEMAP_HOME}/WLS1035/onemap/onemap/webapps
		# fi
	fi
	
	echo "INFO: Finish installing WEB Server"
	return 0
}

# Discription: Install Token Server
# @RET: 0 => Install successfully
#		!0 => Install failed
function instSrvToken(){
	# first, check whether install MoniterAgent
	chkServer "$MA"
	RET=$?
	if test $RET -eq 2 
	then
		instSrvMonitor
	fi
	
	# check whether installed web container
	chkServer "$ONEMAP"
	RET=$?
	if test $RET -eq 0	# the service is running then terminal it first
	then
		service $ONEMAP stop
	elif test $RET -eq 2 
	then 
		# cp -r OneMap/WLS1035 $ONEMAP_HOME
		# if ( ! [ -d ${ONEMAP_HOME}/WLS1035/onemap/onemap/webapps ] );then 
			# mkdir -p $ONEMAP_HOME/WLS1035/onemap/onemap/webapps
		# fi
		
		# cp -r OneMap/webapps/RemoteTokenServer  ${ONEMAP_HOME}/WLS1035/onemap/onemap/webapps
		
		regService "$ONEMAP"
		RET=$?
		if [ $RET -ne 0 ];then
			echo "ERROR: Failed to install onemap Server"
			return 1
		fi	
	# else
		# cp -r OneMap/webapps/RemoteTokenServer  ${ONEMAP_HOME}/WLS1035/onemap/onemap/webapps
	fi
	
	echo "INFO: Finish installing Token Server"
	return 0
}

# Discription: Install Message Server
# @RET: 0 => Install successfully
#		!0 => Install failed
function instSrvMsg(){
	# first, check whether install MoniterAgent
	chkServer "$MA"
	RET=$?
	if test $RET -eq 2 
	then
		instSrvMonitor
	fi

	# copy files
	# if ( ! [ -d ${ONEMAP_HOME}/services ] );then
		# mkdir ${ONEMAP_HOME}/services
	# fi
	# cp -r OneMap/services/activemq5.4.1 $ONEMAP_HOME/services/
	
	# registration the service
	regService "$MQ"
	RET=$?
	if [ $RET -ne 0 ];then
		echo "ERROR: Failed to install ActiveMQ Server"
		return 1
	fi
	
	echo "INFO: Install ActiveMQ Server successfully"
	
	return 0
}

# Discription: Check the OneMap account and create it
# RET: 0 => The OneMap is existed or create it successfully
#     !0 => Create it failed
function chkOMAccount() {
  checkAccount ${OM_ACCOUNT}
  RET=$?
  if test $RET -ne 0
  then
    echo "WARN: The OneMap account isn't existed and then create it!" 

    # add the esri group
    groupadd ${ESRI_GROUP}  
    useradd -g ${ESRI_GROUP} -G root -d /home/${OM_ACCOUNT} ${OM_ACCOUNT}
    echo ${OM_PWD}|passwd --stdin ${OM_ACCOUNT}

    if test $? -ne 0
    then
      echo "ERROR: Create OneMap account failed and return $?"
      return 1
    fi    
  fi

  return 0
}

# Discription: Search the sub string
# RET: 0 => Success
#      1 => Failed
function strstr() {
    declare -i i n2=${#2} n1=${#1}-n2  
    #echo $i $n1 $n2  
    for ((i=0; i<n1; ++i)){  
        if [ "${1:i:n2}" == "$2" ]; then  
            return 0  
        fi  
    }  
    return 1  
}

###############################################################################
# The main control flow
###############################################################################
# Make sure this is being run as root.
tmp=`id | cut -f2 -d\( | cut -f1 -d\)`
if [ "$tmp" != "root" ]
then
   echo " "
   echo "*** ERROR: This script must be run as root."
   echo " "
   exit 1
fi

printIntro

# add the params.conf file
#source params.conf

CONTAINER_HOME=${ONEMAP_HOME}"/"${CONTAINER_NAME}${CONTAINER_VERSION}	# Web server container home
echo "CONTAINER_HOME=$CONATAINER_HOME"

options=$*
echo "Input $# parameters: $options"
# Parse the input parameters and check
if test $# -lt 1 
then
	echo "WARN: Input parameters error."
    printUsage
    exit 1
fi

# check the OneMap account
chkOMAccount
if test $? -ne 0
then
  echo "ERROR: Check OneMap account($OM_ACCOUNT) failed and return!"
  exit 1
fi

for opt in $options;
do
	case "$opt" in
		agent)
			instSrvMonitor
			RET=$?
			if test $RET -eq 0
			then
				chmod 755 -R ${ONEMAP_HOME}
				chown -R ${OM_ACCOUNT}:root ${ONEMAP_HOME}
			fi
			;;
		db)
			instSrvDB
			if test $RET -eq 0
			then
				chmod 755 -R ${ONEMAP_HOME}
				chown -R ${OM_ACCOUNT}:root ${ONEMAP_HOME}
			fi
			;;
		gis)
			instSrvGIS
			if test $RET -eq 0
			then
				chmod 755 -R ${ONEMAP_HOME}
				chown -R ${OM_ACCOUNT}:root ${ONEMAP_HOME}
			fi
			;;
		main)
			instSrvMain
			if test $RET -eq 0
			then
				chmod 755 -R ${ONEMAP_HOME}
				chown -R ${OM_ACCOUNT}:root ${ONEMAP_HOME}
			fi
			;;
		web)
			instSrvWEB
			if test $RET -eq 0
			then
				chmod 755 -R ${ONEMAP_HOME}
				chown -R ${OM_ACCOUNT}:root ${ONEMAP_HOME}
			fi
			;;
		token)
			instSrvToken
			if test $RET -eq 0
			then
				chmod 755 -R ${ONEMAP_HOME}
				chown -R ${OM_ACCOUNT}:root ${ONEMAP_HOME}
			fi
			;;
		msg)
			instSrvMsg
			if test $RET -eq 0
			then
				chmod 755 -R ${ONEMAP_HOME}
				chown -R ${OM_ACCOUNT}:root ${ONEMAP_HOME}
			fi
			;;
		all)
			instSrvMonitor
			if test $RET -ne 0
			then
				echo "Install MonitorAgent Server failed!"
				exit
			fi	
			
			instSrvDB
			if test $RET -ne 0
			then
				echo "Install Database Server failed!"
				exit
			fi	
			
			instSrvGIS
			if test $RET -ne 0
			then
				echo "Install GIS Server failed!"
				exit
			fi			
			
			instSrvMain
			if test $RET -ne 0
			then
				echo "Install Maintenance Server failed!"
				exit
			fi		
			
			instSrvWEB
			if test $RET -ne 0
			then
				echo "Install Web Server failed!"
				exit
			fi			
			
			instSrvToken
			if test $RET -ne 0
			then
				echo "Install Token Server failed!"
				exit
			fi			
			
			instSrvMsg
			if test $RET -ne 0
			then
				echo "Install Message Server failed!"
				exit
			fi
			
			chmod 755 -R ${ONEMAP_HOME}
			chown -R ${OM_ACCOUNT}:root ${ONEMAP_HOME}
			;;
	esac
done

# delete the params.conf file and itself
#rm -f $0
#rm -f params.conf
