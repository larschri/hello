# PCAP

Create a .pcap (packet capture) file with known content, and print the packets
from one tcp connection as json.

Note: the traffic is captured using tshark running with sudo privileges

## pcap file

The pcap file will contain tcp packets from two simultaneous tcp connections
(HTTP requests).
