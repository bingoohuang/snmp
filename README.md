# snmp

snmp package and tools in golang.

## examples

uses the SNMP GET/BULKWALK request to query for information on a network entity

```sh
# snmp -t bmon@192.168.12.13 1.3.6.1.4.1.43353.1.1.1.0 1.3.6.1.4.1.43353.1.1.2.0 .1.3.6.1.4.1.43353.1.1.3.0
[ get][bmon@192.168.12.13][0] .1.3.6.1.4.1.43353.1.1.1.0 = number: 0
[ get][bmon@192.168.12.13][1] .1.3.6.1.4.1.43353.1.1.2.0 = string: SVSJm#2018/05/17-2023/05/16;SVSQm#2018/05/17-2023/05/16;156b1ba46762d0be#2018/05/17-2023/05/16
[ get][bmon@192.168.12.13][2] .1.3.6.1.4.1.43353.1.1.3.0 = string: 2021/1/14 11:9:20
[walk][bmon@192.168.12.13][0] .1.3.6.1.4.1.43353.1.1.1.0 = number: 0
[walk][bmon@192.168.12.13][0] .1.3.6.1.4.1.43353.1.1.2.0 = string: SVSJm#2018/05/17-2023/05/16;SVSQm#2018/05/17-2023/05/16;156b1ba46762d0be#2018/05/17-2023/05/16
[walk][bmon@192.168.12.13][0] .1.3.6.1.4.1.43353.1.1.3.0 = string: 2021/1/14 11:9:20
```

start snmp trap server:

```sh
# snmp -trap :9162                                                                                     
2021/01/14 11:57:32 got trapdata from 127.0.0.1
[ trap][127.0.0.1:49784][0] .1.3.6.1.2.1.1.3.0 = number: 87652648
[ trap][127.0.0.1:49784][1] .1.3.6.1.6.3.1.1.4.1.0 = number: 0
[ trap][127.0.0.1:49784][2] .1.3.6.1.4.1.8072.2.3.2.1 = number: 123456
2021/01/14 12:01:21 got trapdata from 127.0.0.1
[ trap][127.0.0.1:56368][0] .1.3.6.1.2.1.1.3.0 = number: 87675610
[ trap][127.0.0.1:56368][1] .1.3.6.1.6.3.1.1.4.1.0 = number: 0
[ trap][127.0.0.1:56368][2] .1.3.6.1.4.1.43353.1.1.2.0 = string: bingoohuang
```

Send A Test Trap:

```sh
# snmptrap -v 2c -c public localhost:9162 '' 1.3.6.1.4.1.8072.2.3.0.1 1.3.6.1.4.1.43353.1.1.2.0  s bingoohuang
# snmptrap -v 2c -c public localhost:9162 '' 1.3.6.1.4.1.8072.2.3.0.1 1.3.6.1.4.1.8072.2.3.2.1 i 123456
```

## resources

1. [SNMP Trap - How To Send A Test Trap](https://support.nagios.com/kb/article.php?id=493)

SNMP v2 Trap

- Command form: `snmptrap -v <snmp_version> -c <community> <destination_host> <uptime> <OID_or_MIB> <object> <value_type> <value>`
- Using MIB: `snmptrap -v2c -c public localhost '' NET-SNMP-EXAMPLES-MIB::netSnmpExampleHeartbeatNotification netSnmpExampleHeartbeatRate i 123456`
- Shortening MIB: `snmptrap -v2c -c public localhost '' netSnmpExampleHeartbeatNotification netSnmpExampleHeartbeatRate i 123456`
- Using OID: `snmptrap -v 2c -c public localhost '' 1.3.6.1.4.1.8072.2.3.0.1 1.3.6.1.4.1.8072.2.3.2.1 i 123456`

The commands above required the following settings in /etc/snmp/snmptrapd.conf

    disableAuthorization yes
    traphandle default /usr/sbin/snmptthandler  
