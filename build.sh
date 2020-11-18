#go build -ldflags "-s -w" -o ~/iserver main.go
#go build -ldflags "-s -w" -o ./iserver main.go

#https://studygolang.com/articles/10763
#go build -ldflags '-w -s'

#$ upx etcd-cli
## 此处省略压缩时的打印...
#$ ls -lh
#-rwxr-xr-x  1 gangan  staff   897K Aug 18 00:49 etcd-cli
#-rw-r--r--  1 gangan  staff   456B Aug 18 00:34 main.go