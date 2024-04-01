json: log.pcap
	 @tshark -r log.pcap -T json|jq '.[]._source.layers|select(.tcp["tcp.stream"]=="0")'

log.pcap:
	(go run ./... :8080 &) 
	(sudo timeout -s 2 2s ngrep -O log.pcap -d lo port 8080 | wc  &) 
	sleep 1.1
	(curl http://localhost:8080/ &)
	curl http://localhost:8080/
	curl http://localhost:8080/exit
	sleep 1.1

