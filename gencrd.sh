rm -rf ./crd/generated
go mod vendor

# 第二行需填写项目完整的 gopath
bash vendor/k8s.io/code-generator/generate-groups.sh all \
github.com/schwarzeni/k8senhance/crd/generated github.com/schwarzeni/k8senhance/crd/api \
cloudedgeservice:v1alpha1 \
--go-header-file hack/boilerplate.go.txt \
--output-base ../../../

rm -rf vendor
