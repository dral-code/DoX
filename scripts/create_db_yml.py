import os

filename = "files/vars.yml"
if os.path.exists(filename):
    os.remove(filename)

file = open (filename, "a")
file.write("records:\n")

for domain in range(0,251):
    for host in range(1,251):
        #print ("    host" + str(host) + ".domain" + str(domain) + ": {forward: 10.0." + str(domain) + "." + str(host) + ", type: A, last: " + str(host) + "." + str(domain) + ", rev: 0.10.in-addr.arpa.}\n")
        file.write ("    host" + str(host) + ".domain" + str(domain) + ": {forward: 10.0." + str(domain) + "." + str(host) + ", type: A, last: " + str(host) + "." + str(domain) + ", rev: 0.10.in-addr.arpa.}\n")

file.close()            