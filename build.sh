 set -e
 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o k8senhance

 # 用于网络和调度测试
# scp k8senhance root@10.211.55.52:~/
# scp k8senhance root@10.211.55.65:~/
# scp k8senhance root@10.211.55.66:~/

# 用于镜像加速测试
scp k8senhance root@10.211.55.67:~/
scp k8senhance root@10.211.55.68:~/
scp k8senhance root@10.211.55.69:~/
scp k8senhance root@10.211.55.70:~/
