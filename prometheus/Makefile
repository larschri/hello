

start:
	docker run \
		--volume ${PWD}/prometheus.yml:/etc/prometheus/prometheus.yml \
		--detach \
		--network=host \
		--rm \
		--name helloprometheus \
		prom/prometheus

rungo1:
	go run ./... :9042

rungo2:
	go run ./... :9043

stop:
	docker stop helloprometheus
