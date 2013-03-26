#!/bin/bash

####################################################################################################
# DISP: Uninstall bash script for OneMap Linux version
# AUTH: Michael.Cho
# TIME: 2012/12/30
# EMAIL: zuow11@gmail.com
####################################################################################################

####################################################################################################
# consts and virables
####################################################################################################
MA="monitoragent"							# MonitorAgent service in onemap
MQ="activemq"								# Activemq service in onemap
H2="h2memdb"								# H2memdb service in onemap
ONEMAP="onemap"								# Onemap service 

# @Discription: unregistration a service
# @Arguments: unregService SERVICE_NAME
# @RET: 0  => unregistration service successfully
#       !0 => unregistration failed
function unregService(){
    if [ $# -eq 0 ];then
	echo "ERROR: Please input the unregistration service name"
	return 1
    elif [ $# -ne 1 ];then
	echo "ERROR: Please input only one service name per time"
	return 1
    fi

    svcName=$1
	
    # first check the status of service 
    chkService $svcName
    local RET=$?
    if test $RET -eq 0  # the service is running
    then
        # shut down the service
        /etc/init.d/${svcName} stop

        sleep 5000

        if test -f "/var/run/$svcName.pid"
        then
            PID=`cat /var/run/$svcName.pid`
            if test -n $PID
            then
                 kill -9 $PID
            fi

            rm -f /var/run/$svcName*
        fi
    fi

    # clear the service
    if test -f "/etc/init.d/$svcName"
    then
        rm -f /etc/init.d/$svcName
    fi

    # delete the null pid file
    if test -f "/var/run/${svcName}.pid"
    then
        rm -f /var/run/${svcName}.pid
    fi

    return 0
}

# Discription: Check a service running or not
# Arguments: chkService SERVER_NAME
# @RET: 0 => the service is running
# 		1 => the service isn't running
#		2 => the service isn't installed
#		3 => others
function chkService(){
    svcName=$1
	
    if test -f "/etc/init.d/$svcName"	
    then
        if test -f "/var/run/${svcName}.pid"
	then
            if test -z "$(cat /var/run/$svcName)"
	    then 
		echo "INFO: The pid file of Service $svcName is null"
		return 3
	    fi
			
	    if test ps -p $(cat /var/run/$svcName) >/dev/null
	    then
		echo "INFO: Service $svcName(pid: cat /var/run/${svcName}.pid)"
		return 0
	    else
		echo "INFO: Service $svcName has terminated unexpectedly"
		return 1
	    fi
        else
	    echo "INFO: Service $svcName isn't running"
	    return 1
    	fi
    else
	echo "INFO: Service $svcName hasn't installed"
	return 2
    fi
}

# Discription: Delete the db data
# @RET: 0  => success
#       !0 => failed
function delDB(){
    # first check the uninstall scrip whether existed
    if test -f "${ONEMAP_HOME}/db/uninstall.sql"
    then
        ## please backup the database first

        # exec the uninstall sql script
	    su - ${ORCL_ACCOUNT} -c "/bin/sh -c \"echo exit | $ORACLE_HOME/bin/sqlplus ${ORACLE_SYSTEM_ACCOUNT}/${ORACLE_SYSTEM_PWD} as sysdba @$ONEMAP_HOME/db/uninstall.sql \"" 
        if test $? -ne 0
        then
            echo "ERROR: Exec the unistall db script failed"
            return 1
        fi
        
        # get the oracle base 
        ORACLE_BASE=`grep "<PROPERTY NAME=\"ORACLE_BASE\"" $ORACLE_HOME/inventory/ContentsXML/oraclehomeproperties.xml|awk '{print $3}'|awk -F\" '{print $2}'`
        if test -z $ORACLE_BASE
        then
            echo "ERROR: Get ORACLE_BASE failed and please check!"
            return 2
        fi
        ORACLE_DATA=${ORACLE_BASE}"/oradata/"${ORACLE_SID}
        
        # delete the database files
        if test -f "${ORACLE_DATA}/GEO*.dbf"
        then
            su - ${ORCL_ACCOUNT} -c "rm -f ${ORACLE_DATA}/GEO*.dbf"
        fi
    else
	echo "ERROR: There isn\'t uninstall.sql file"
    fi

    return 0
}

#### MAIN
# uninstall the services
unregService "$MA"
RET=$?
if test ${RET} -ne 0
then
    echo "ERROR: Uninstall $MA service failed"
fi
unregService "$H2"
RET=$?
if test ${RET} -ne 0
then
    echo "ERROR: Uninstall $H2 service failed"
fi
unregService "$MQ"
RET=$?
if test ${RET} -ne 0
then
    echo "ERROR: Uninstall $MQ service failed"
fi
unregService "$ONEMAP"
RET=$?
if test ${RET} -ne 0
then
    echo "ERROR: Uninstall $ONEMAP service failed"
    return 1
fi

# delete the db 
delDB
RET=$?
if test ${RET} -ne 0
then
    echo "ERROR: Uninstall db failed"
    return 2
fi

# delete the OneMap package
rm -rf ${ONEMAP_HOME}
RET=$?
if test ${RET} -ne 0
then
    echo "ERROR: Delete the OneMap package failed"
    return 3
fi

echo "Uninstall OneMap successfully"
