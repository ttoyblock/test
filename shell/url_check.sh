#!/bin/bash

SMS='http://192.168.0.72:80/tools/sms.php?source=url_check'

projectServerList=(
'192.168.0.67:20017'
'192.168.0.72:20017'
'192.168.0.65:123123'
)

TO='"13681487657","18380465345"'

for((i=0;i<${#projectServerList[@]};i++))
do
	echo ${projectServerList[$i]}
	#接口探活
	pingInfo=`curl -s ${projectServerList[$i]}/ping`

	pingInfo=${pingInfo//\"/}
	if [[ $pingInfo != "ping" ]];then

		# log
		echo ${projectServerList[$i]}" server ping error "$pingInfo

		# send sms
		msg=${projectServerList[$i]}
		curl -s -H "Content-Type: application/json" -X POST --data '{"msg":"'$msg'端口不可用!","to":['$TO']}' ${SMS}
		# echo $curlinfo
	fi
done