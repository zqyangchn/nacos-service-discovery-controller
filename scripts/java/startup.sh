#!/bin/bash

export SERVER_PORT=8080

java -server -Xms4g -Xmx4g -XX:+UseG1GC -XX:SurvivorRatio=3 -XX:G1ReservePercent=15 -XX:MaxGCPauseMillis=100 -Dspring.profiles.active=sit -Dfile.encoding=UTF-8 -Duser.timezone=GMT+07 -Drocketmq.client.logFileMaxSize=134217728 -Drocketmq.client.logFileMaxIndex=2 -Dxxl.job.admin.addresses=https://sit-th-flc-xxl.flashfin.com -Dnacos.config.server-addr=https://dev-nacos.flashfin.com -jar flm-flc-cg-mservice.jar
