#! /bin/bash
export binaryName=zserver
export imageName=zserver:latest
export containerName=zserver

function build() {
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -ldflags "-s -w" -o $binaryName cmd/main.go
    echo "build success"
}

function iosbuild() {
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $binaryName
    echo "ios build success"
}

function dockerBuild() {
    docker build --build-arg  -t $imageName .
    echo "dev docker build success"
}

function dockerStop() {
    docker stop $containerName && docker rm -f $containerName
    echo "docker stop success"
}

function dockerRun() {
    docker run --restart=always --net=host --name $containerName  -idt -p 8000:8000  -v /etc/hosts:/etc/hosts  -v /var/log/:/var/log/ $imageName
    echo "docker run success"
}

function dockerClean() {
    docker images | grep none | awk '{print $3 }' | xargs docker rmi
}

function dockerDeploy() {
     build
     dockerBuild
     dockerStop
     dockerRun
}

function tailf() {
   tail -f $logfile
}

function help() {
    echo "$0 start|stop|restart"
}

if [ "$1" == "" ]; then
    help
elif [ "$1" == "build" ];then
    build
elif [ "$1" == "iosbuild" ];then
    iosbuild
elif [ "$1" == "dockerRun" ];then
    dockerRun
elif [ "$1" == "dockerStop" ];then
    dockerStop
elif [ "$1" == "dockerDeploy" ];then
    dockerDeploy
elif [ "$1" == "dockerClean" ];then
    dockerClean
else
    help
fi

