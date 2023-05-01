import os

filename = "../shared_files/list_of_dns_entry.csv"
if os.path.exists(filename):
    os.remove(filename)

file = open (filename, "a")
file.write("dnsEntryName;dnsEntryIP\n")

for domain in range(0,251):
    for host in range(1,251):
        #print ("host" + str(host) + ".domain" + str(domain) + ".rsx218-dox.cnam.fr;172.16." + str(domain) + "." + str(host))
        file.write ("host" + str(host) + ".domain" + str(domain) + ".rsx218-dox.cnam.fr;172.16." + str(domain) + "." + str(host) +"\n")

file.close()            