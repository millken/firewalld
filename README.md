# firewalld

dbi #dynamic blacklist ip
dbn # dynamic blacklist net
sbi #static blacklist ip
sbn #static blacklist net
dwi #dynamic whitelist ip
dwn #dynamic whitelist net
swi #static whitelist ip
swn #static whitelist net

ipset create -exist dbi hash:ip hashsize 4096 maxelem 1048576
ipset create -exist dbn hash:net hashsize 4096 maxelem 1048576
ipset create -exist sbi hash:ip hashsize 4096 maxelem 1048576
ipset create -exist sbn hash:net hashsize 4096 maxelem 1048576
ipset create -exist dwi hash:ip hashsize 4096 maxelem 1048576
ipset create -exist dwn hash:net hashsize 4096 maxelem 1048576
ipset create -exist swi hash:ip hashsize 4096 maxelem 1048576
ipset create -exist swn hash:net hashsize 4096 maxelem 1048576
