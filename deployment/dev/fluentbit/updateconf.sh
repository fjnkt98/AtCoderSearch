set -e

microk8s kubectl create configmap fluent-bit-config --from-file=fluent-bit.conf --dry-run=client --output=yaml | microk8s kubectl apply -f -
microk8s kubectl create configmap parsers-config --from-file=parsers.conf --dry-run=client --output=yaml | microk8s kubectl apply -f -
microk8s kubectl create configmap backend-config --from-file=backend.conf --dry-run=client --output=yaml | microk8s kubectl apply -f -
microk8s kubectl create configmap batch-config --from-file=batch.conf --dry-run=client --output=yaml | microk8s kubectl apply -f -
