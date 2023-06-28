import sys
import pandas as pd
import math
import datetime

# no check of file path quick and dirty
filepath = sys.argv[1]

csvfile = pd.read_csv(filepath, sep=',')

line_count = 0
req_total = float(0)
for index, line in csvfile.iterrows():
    req_count = float(line['Incoming'])
    if int(round(req_count,0)) < 4:
        continue
    line_count += 1
    req_total += req_count
    print("request num : " + str(line_count) + " sys " + str(line['Incoming']) + " total " + str(req_total))
#    if line_count == 5:
#        break

usage_total_round = round(req_total,2)
print("nb de ligne : " + str(line_count) + "=> temps : " + str(datetime.timedelta(seconds=line_count*15)) + " for " + str(req_total) + " requetes moyenne " + str(round(req_total/line_count,2)) + "req/s")