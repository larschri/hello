# PCAP

Create a .pcap (packet capture) file with known content, and print the packets
from one tcp connection as json.

## pcap file

The pcap file will contain tcp packets from three tcp connections (HTTP
requests). 

## Dependencies

tshark, ngrep, curl, jq, go, timeout
