.PHONY: create-configmap
create-configmap:
	@microk8s kubectl create configmap db-migrations --from-file=../../db/migrations --dry-run=client -o=yaml | microk8s kubectl apply -f -
	@microk8s kubectl create configmap fluent-bit-config --from-file=daemonset/fluent-bit.conf --dry-run=client --output=yaml | microk8s kubectl apply -f -
	@microk8s kubectl create configmap parsers-config --from-file=daemonset/parsers.conf --dry-run=client --output=yaml | microk8s kubectl apply -f -
