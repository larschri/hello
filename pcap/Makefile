json: log.pcap
	 tshark -r log.pcap -T json|jq '.[]._source.layers|select(.tcp["tcp.stream"]=="0")|select(.tcp["tcp.payload"] != null)'

log.pcap:
	sudo touch log.pcap
	(sudo tshark -Q -w log.pcap -i lo -F libpcap -a duration:2 -f 'port 8080' &)
	sleep 0.5
	go run ./...
	sleep 2
