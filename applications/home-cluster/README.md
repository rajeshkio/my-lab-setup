Running application at home lab

Challenges:

1 - No public Access

2 - Issues with SSL

I solved them by below process:

 - Subscribed to no-ip dynamic dns service and bought a domain. This will make sure domain is always updated with my latest public IP.

 - Added CNAME records for route53 subdomains of my public domain. CNAME record was to redirect traffic from my public domain to dynamic dns name. E.g. grafana.rajesh-kumar.in --> mydynamicdns.ddns.net

 - I already had certificates for my public domain. Just needed to create secret with those certificate in application specific namespaces/

 - Disabled https access through WAN and LAN from router security --> access control settings. After this 443 port forward is working

 - Additionally I added domains like grafana.rajesh-kumar.in into router DNS so that I can resolve them locally and get the access. To access it from the outside, I had to use domain:port.
