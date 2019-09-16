import sys
import os
f=open(sys.argv[1])
out_str=""
out_ip=""
for line in f:
    line = line.strip()
    if not line.startswith("fab deploy.deploy_go_server"):
        continue
    terms=line.split(" ")
    print terms[5]
    raw_ip=os.popen("host " + terms[3]).read()
    ips=raw_ip.strip().split(" ")
    print "ip:",ips[3]
    out_str+=",\"" +terms[5] + "\""
    out_ip+=",\"" +ips[3] + "\""
print out_str
print out_ip
