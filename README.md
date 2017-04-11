# firewalld

dbi #dynamic blacklist ip
dbn # dynamic blacklist net
sbi #static blacklist ip
sbn #static blacklist net
dwi #dynamic whitelist ip
dwn #dynamic whitelist net
swi #static whitelist ip
swn #static whitelist net


# first delete iptables rules
ipset destroy dbi
ipset destroy dbn
ipset destroy sbi
ipset destroy sbn
ipset destroy dwi
ipset destroy dwn
ipset destroy swi
ipset destroy swn
ipset create -exist dbi hash:ip hashsize 4096 maxelem 1048576 timeout 86400
ipset create -exist dbn hash:net timeout 86400
ipset create -exist sbi hash:ip hashsize 4096 maxelem 1048576
ipset create -exist sbn hash:net hashsize 2048 maxelem 524288
ipset create -exist dwi hash:ip hashsize 2048 maxelem 524288 timeout 86400
ipset create -exist dwn hash:net timeout 86400
ipset create -exist swi hash:ip hashsize 4096 maxelem 1048576
ipset create -exist swn hash:net hashsize 2048 maxelem 524288
iptables -A INPUT  -m set --match-set dbi src -p TCP -m multiport --dports 80,12377 -j DROP
iptables -A INPUT  -m set --match-set dbn src -p TCP -m multiport --dports 80,12377 -j DROP
iptables -A INPUT  -m set --match-set dbn sbi -p TCP -m multiport --dports 80,12377 -j REJECT
iptables -A INPUT  -m set --match-set dbn sbn -p TCP -m multiport --dports 80,12377 -j REJECT
iptables -I INPUT  -m set --match-set dwi src -p TCP -m multiport --dports 80,12377 -j ACCEPT
iptables -I INPUT  -m set --match-set dwn src -p TCP -m multiport --dports 80,12377 -j ACCEPT
iptables -I INPUT  -m set --match-set swi src -p TCP -m multiport --dports 80,12377 -j ACCEPT
iptables -I INPUT  -m set --match-set swn src -p TCP -m multiport --dports 80,12377 -j ACCEPT

