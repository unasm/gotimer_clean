#########################################################################
# File Name :    make.sh
# Author :       unasm
# mail :         unasm@sina.cn
# Last_Modified: 2016-08-24 21:40:58
#########################################################################
#!/bin/bash
OLDGOPATH="$GOPATH"
echo '$0: '$0

cd `dirname $0` 

export GOPATH=`pwd`
export CGO_ENABLED=0

gofmt -w src

go install security

export GOPATH="$OLDGOPATH"

echo 'ok.'

./bin/security
#go get github.com/beego/bee
#go get github.com/astaxie/beego
