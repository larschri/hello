all: cpu.out
	go tool pprof -http :8081 cpu.out

cpu.out:
	go test -bench=. -cpuprofile=cpu.out .

clean:
	rm cpu.out

