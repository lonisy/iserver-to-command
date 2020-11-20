#go build -ldflags "-s -w" -o ~/bin/iserver main.go
#go build -ldflags "-s -w" -o ./iserver main.go

#https://studygolang.com/articles/10763
#go build -ldflags '-w -s'

# 跨平台编译
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./iserver ./main.go
#CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./iserver ./main.go

#go get -u github.com/mattn/go-sqlite3
#go get -u github.com/mitchellh/go-homedir

#$ upx etcd-cli
## 此处省略压缩时的打印...
#$ ls -lh
#-rwxr-xr-x  1 gangan  staff   897K Aug 18 00:49 etcd-cli
#-rw-r--r--  1 gangan  staff   456B Aug 18 00:34 main.go


function to(){
  echo "func from bashrc."
  aliasCmd=`iserver to $@`
  if [[ "${aliasCmd:0:8}" == "Commands" ]]; then
    echo -e "\033[32m${aliasCmd} \033[0m"
    IFS=$'\n'; arrIN=($aliasCmd); unset IFS;
    for line in "${arrIN[@]}"; do
      ${line:10}
      break
    done
  else
    echo -e "\033[32m${aliasCmd} \033[0m"
  fi
}