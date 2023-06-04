# DoX
RSX218 projet DoX

# ansible
I'm using the following version:
```
ansible [core 2.12.10]
```

# monitoring stack

I used a classic prometheus/grafana monitoring stack as it's easy to configure.

## installation
I used the following tutorial to install the monitoring stack on the host:
* https://grafana.com/docs/grafana/latest/setup-grafana/installation/debian/
* https://www.fosstechnix.com/install-prometheus-and-grafana-on-ubuntu/

## dashboard
I'm using this dashboard for visualization https://grafana.com/grafana/dashboards/1860-node-exporter-full/

# PKI
I used the following website as inspiration for my certificate generation:
https://gquintana.github.io/2020/11/28/Build-your-own-CA-with-Ansible.html

Dont't forget to install the ansible module requirements specified on the website.

# Resolver/Forwarder

## Do53

As we use two solutions (bind9 and DNS Proxy), we need a baseline for each.

### bind version
The bind9-isc package is used as forwarder for Do53. It's working out of the box after the deployment.

### DNS Proxy version
To run DNS Proxy in plain mode:
```
sudo ./dnsproxy -l 192.168.56.5 -u udp://192.168.56.2:453 -v
```

## DoT
The bind9-isc package is used as forwarder for DoT. It's working out of the box after the deployment.

## DoH
The bind9-isc package is used as forwarder for DoH. It's working out of the box after the deployment.

## DoQ
DNS Proxy https://github.com/AdguardTeam/dnsproxy is used as intermediate resolver as there is no other alternative.
Once the server is deployed, you need to run the following command to start it:
```
sudo ./dnsproxy -l 192.168.56.5 --quic-port=853 --tls-crt=/etc/pki/tls/private/rsx218-dox.crt --tls-key=/etc/pki/tls/private/rsx218-dox.key -u udp://192.168.56.2:453 -p 0 -v'
```
