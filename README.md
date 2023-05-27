# DoX
RSX218 projet DoX

# ansible
I'm using the following version :
'ansible [core 2.12.10]'

# monitoring stack

I used a classic prometheus/grafana monitoring stack as it's easy to configure.

## installation
I used the following tutorial to install the monitoring stack on the host :
* https://grafana.com/docs/grafana/latest/setup-grafana/installation/debian/
* https://www.fosstechnix.com/install-prometheus-and-grafana-on-ubuntu/

## dashboard
I'm using this dashboard for visualization https://grafana.com/grafana/dashboards/1860-node-exporter-full/

# PKI
I used the following website as inspiration for my certificate generation :
https://gquintana.github.io/2020/11/28/Build-your-own-CA-with-Ansible.html

Dont't forget to install the ansible module requirements specified on the website.