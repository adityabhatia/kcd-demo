# modularization-demo

### Two cluster mode with Gardener shoots
>The following steps are applicable each individual module.

#### Module-operator bundling
```shell
## Module Operator
# If template.yaml is already prepared - SKIP this section

# Ensure gcloud is authenticated
export MODULE_REGISTRY=europe-west3-docker.pkg.dev/sap-kyma-jellyfish-dev/operator-demo  
export IMG_REGISTRY=$MODULE_REGISTRY/operator-images
export MODULE_CREDENTIALS=oauth2accesstoken:$(gcloud auth print-access-token --impersonate-service-account operator-demo-sa@sap-kyma-jellyfish-dev.iam.gserviceaccount.com)
make module-image
make module-build
```

#### Repo links to run `make` commands below:
* [lifecycle-manager](https://github.com/kyma-project/lifecycle-manager/tree/main/operator)
* [module-manager](https://github.com/kyma-project/module-manager/tree/main/operator)

```shell
## KCP cluster

# Deploy module-manager 
make deploy IMG=eu.gcr.io/kyma-project/module-manager:latest
# Deploy lifecycle-manager 
make deploy IMG=eu.gcr.io/kyma-project/lifecycle-manager:latest

# Secret to SKR cluster kubeconfig
kubectl create ns kyma-system
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: kyma-sample #change with your kyma name
  namespace: kyma-system
  labels:
    "operator.kyma-project.io/managed-by": "lifecycle-manager"
    "operator.kyma-project.io/kyma-name": "kyma-sample"
type: Opaque
data:
  config: $(cat <path-to-kubeconfig-file> | sed 's/---//g' | base64)
EOF

# Install module template
kubectl apply -f ./template.yaml

# Generate kyma resource to start installation
/bin/bash ./hack/gen-kyma.sh
kubectl apply -f kyma.yaml 
```


