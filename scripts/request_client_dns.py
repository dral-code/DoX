import subprocess
import random
import datetime
import socket

n = 250
ihost = [i for i in range (1, n+1)]
random.shuffle(ihost)
idomain = [i for i in range (0, n+1)]
random.shuffle(idomain)
print("client,reqID,url,timestamp")
hostname = socket.gethostname()

line_count = 0
for d in idomain:
    for h in ihost:
        dev_null = " > /dev/null 2>&1"
        #dev_null = ""
        url = "host" + str(h) + ".domain" + str(d) + ".rsx218-dox.cnam.fr"
        cmd_str = "dnslookup host" + url + " 192.168.56.2:453" + dev_null + " &"
        subprocess.run(cmd_str, shell=True)
        line_count += 1
        print(hostname + "," + str(line_count) + "," + url + "," + str(datetime.datetime.now()))
#        if line_count == 5:
#            break