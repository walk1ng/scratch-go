
# 代码生成的工作目录，也就是我们的项目路径
# ROOT_PACKAGE="github.com/walk1ng/my-k8s-controller-custom-resource"

MODULE_NAME="my-k8s-controller-custom-resource"

# API Group
CUSTOM_RESOURCE_NAME="samplecrd"
# API Version
CUSTOM_RESOURCE_VERSION="v1"

# 安装k8s.io/code-generator
# go get k8s.io/code-generator/...
# cd $GOPATH/src/k8s.io/code-generator

# 执行代码自动生成，其中pkg/client是生成目标目录，pkg/apis是类型定义目录
$GOPATH/pkg/mod/k8s.io/code-generator@v0.19.2/generate-groups.sh all "$MODULE_NAME/pkg/client" "$MODULE_NAME/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"