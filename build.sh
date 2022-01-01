 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o k8senhance
 scp k8senhance root@10.211.55.52:~/
 scp k8senhance root@10.211.55.65:~/
 scp k8senhance root@10.211.55.66:~/
