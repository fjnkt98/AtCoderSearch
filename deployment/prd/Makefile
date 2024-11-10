.PHONY: create-configmap
create-configmap:
	@kubectl create configmap db-migrations --from-file=../../db/migrations --dry-run=client -o=yaml | kubectl apply -f -
	@kubectl create configmap fluent-bit-config --from-file=daemonset/fluent-bit.conf --dry-run=client --output=yaml | kubectl apply -f -
	@kubectl create configmap parsers-config --from-file=daemonset/parsers.conf --dry-run=client --output=yaml | kubectl apply -f -
