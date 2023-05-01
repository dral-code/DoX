import os
import shutil

filename = "main.yml"
filepath = "../shared_files/" + filename
if os.path.exists(filepath):
    os.remove(filepath)

file = open (filepath, "a")
file.write("---\n# vars file for Bind setup\ndomain: rsx218-dox.cnam.fr\n\nrev_domain: 16.172.in-addr.arpa\n\nrecords:\n")

for domain in range(0,251):
    for host in range(1,251):
        #print ("    host" + str(host) + ".domain" + str(domain) + ": {forward: 172.16." + str(domain) + "." + str(host) + ", type: A, last: " + str(host) + "." + str(domain) + ", rev: 16.172.in-addr.arpa.}\n")
        file.write ("    host" + str(host) + ".domain" + str(domain) + ": {forward: 172.16." + str(domain) + "." + str(host) + ", type: A, last: " + str(host) + "." + str(domain) + ", rev: 16.172.in-addr.arpa.}\n")

file.close()
shutil.copyfile(filepath, "../ansible/roles/dns-bind9/vars/" + filename)
