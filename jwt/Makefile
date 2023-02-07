all: public-key.pem private-key.pem
	go run ./... private-key.pem public-key.pem

public-key.pem: private-key.pem
	openssl rsa -in private-key.pem -pubout -out public-key.pem

private-key.pem:
	openssl genrsa -out private-key.pem 2048



