
s-ui/
s-ui/sui
s-ui/s-ui.service
s-ui/s-ui.sh
Migration... 
Database not found
Install/update finished! For security it's recommended to modify panel settings 
Do you want to continue with the modification [y/n]? :cancel...
this is a fresh installation,will generate random login info for security concerns:
###############################################
username:Eq1d0KYC
password:OVk7JH2r
###############################################
if you forgot your login info,you can type s-ui for configuration menu
reset admin credentials success
First admin credentials:
        Username:        Eq1d0KYC
        Password:        OVk7JH2r
Created symlink /etc/systemd/system/multi-user.target.wants/s-ui.service → /etc/systemd/system/s-ui.service.
s-ui v1.2.2 installation finished, it is up and running now...
You may access the Panel with following URL(s):
Local address:
http://104.225.239.31:2095/app/
http://169.254.254.254:2095/app/

Global address:
http://104.225.239.31:2095/app/



[root@fancy-byte-3 ~]# for d in www.oracle.com cua-chat-ui.tesla.com www.aws.com www.icloud.com download.amd.com ts3.tc.mm.bing.net r.bing.com c.s-microsoft.com xp.apple.com polyfill-fastly.io ; do t1=$(date +%s%3N); timeout 1 openssl s_client -connect $d:443 -servername $d </dev/null &>/dev/null && t2=$(date +%s%3N) && echo "$d: $((t2 - t1)) ms" || echo "$d: timeout"; done
www.oracle.com: 77 ms
cua-chat-ui.tesla.com: 72 ms
www.aws.com: 91 ms
www.icloud.com: 60 ms
download.amd.com: 79 ms
ts3.tc.mm.bing.net: 25 ms
r.bing.com: 92 ms
c.s-microsoft.com: 59 ms
xp.apple.com: 48 ms
polyfill-fastly.io: 27 ms


static.cloud.coveo.com: 72 ms
amd.com: 35 ms
aws.com: 61 ms
j.6sc.co: 80 ms
www.sony.com: 92 ms
cdn.userway.org: 42 ms
apps.mzstatic.com: 28 ms
digitalassets.tesla.com: 28 ms
gray-wowt-prod.gtv-cdn.com: 140 ms
res-1.cdn.office.net: 96 ms


ssh -L 2095:127.0.0.1:2095 root@104.225.239.31
fireman@123456