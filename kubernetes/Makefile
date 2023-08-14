
KIND=~/go/bin/kind
TAG=latest

run: loadimage
	kubectl --context=kind-kind create -f deployment.yaml

installkind:
	go install sigs.k8s.io/kind@latest

startkind: installkind
	$(KIND) get clusters | grep -q kind || $(KIND) create cluster --config=cluster.yaml

build:
	go build main.go

buildimage: build
	docker build . -t my-image:$(TAG)

loadimage: buildimage startkind
	$(KIND) load docker-image my-image:$(TAG)

clean:
	$(KIND) delete cluster
	rm main
