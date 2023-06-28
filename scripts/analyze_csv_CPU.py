import sys
import pandas as pd
import math
import datetime

# no check of file path quick and dirty
filepath = sys.argv[1]

csvfile = pd.read_csv(filepath, sep=',')

line_count = 0
usage_total = float(0)
for index, line in csvfile.iterrows():
    
    usage_System = float(line['Busy System'].replace('%',''))
    usage_User = float(line['Busy User'].replace('%',''))
    usage = round(math.fsum([usage_System,usage_User]),2)
    round_usage = int(round(usage,0))
    if int(round(usage,0)) < 5:
        continue
    line_count += 1
    usage_total += usage
    print("request num : " + str(line_count) + " sys " + line['Busy System'] + " user " + line['Busy User'] + " round : " + str(round_usage) + " usage " + str(usage))
#    if line_count == 5:
#        break

usage_total_round = round(usage_total,2)
print("nb de ligne : " + str(line_count) + "=> temps : " + str(datetime.timedelta(seconds=line_count*15)) + " moyenne " + str(round(usage_total_round/line_count,2)) + "%")