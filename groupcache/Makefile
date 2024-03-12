run:
	(timeout 2 go run main.go -addr=:7310 -pool=http://127.0.0.1:7310,http://127.0.0.1:7311,http://127.0.0.1:7312) &
	(timeout 2 go run main.go -addr=:7311 -pool=http://127.0.0.1:7311,http://127.0.0.1:7310,http://127.0.0.1:7312) &
	(timeout 2 go run main.go -addr=:7312 -pool=http://127.0.0.1:7312,http://127.0.0.1:7310,http://127.0.0.1:7311;) & 
	sleep 1
	curl localhost:7310/color?name=red
	curl localhost:7310/color?name=green
	curl localhost:7310/color?name=blue
	curl localhost:7311/color?name=red
	curl localhost:7311/color?name=green
	curl localhost:7311/color?name=blue
	sleep 1
