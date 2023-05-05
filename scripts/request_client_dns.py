import subprocess
import pandas 

#csvfile = pandas.read_csv('../shared_files/list_of_dns_entry.csv', sep=';')
csvfile = pandas.read_csv('list_of_dns_entry.csv', sep=';')
csvfile = csvfile.sample(frac=1).reset_index(drop=True)
#print( csvfile)

line_count = 0
for index, line in csvfile.iterrows():
    dev_null = " &> /dev/null"
    #dev_null = ""
    cmd_str = "dnslookup " + line["dnsEntryName"] + " 192.168.56.2:453" + dev_null
    # ip: " + line['dnsEntryIP'])
    subprocess.run(cmd_str, shell=True)
    line_count += 1
    print("request num : " + str(line_count) + " for " + line["dnsEntryName"])
#    if line_count == 5:
#        break