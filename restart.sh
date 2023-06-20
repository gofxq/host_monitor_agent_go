BIN_NAME='host_monitor_agent_go'
DEPLOY_SERVER='cd'
DEPLOY_PATH="run/${BIN_NAME}"

rm -rf ./build
mkdir -p ./build

#cp -r config.yml ./build
CGO_ENABLED=0 \
 GOOS=linux \
 GOARCH=amd64 go build -ldflags="-s -w" -o ./build/${BIN_NAME}_new main.go && upx -1 ./build/${BIN_NAME}_new

#deploy

ssh $DEPLOY_SERVER "mkdir -p ${DEPLOY_PATH}"
scp -r build/* ${DEPLOY_SERVER}:${DEPLOY_PATH}
ssh $DEPLOY_SERVER "mv -f ${DEPLOY_PATH}/${BIN_NAME}_new ${DEPLOY_PATH}/${BIN_NAME}"
ssh $DEPLOY_SERVER "pidof ${BIN_NAME}|xargs kill"
ssh $DEPLOY_SERVER "cd ${DEPLOY_PATH}; nohup ./${BIN_NAME} -addr host-cd.gofxq.com:8008 > ./nohup.out 2>&1   &"
ssh $DEPLOY_SERVER "ps aux|grep ${BIN_NAME}"
ssh $DEPLOY_SERVER "tail -f ${DEPLOY_PATH}/nohup.out"
