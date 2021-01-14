# snmp

snmp package and tools in golang.

## examples

uses the SNMP GET/BULKWALK request to query for information on a network entity

```sh
$ snmp -t bmon@192.168.12.13 1.3.6.1.4.1.43353.1.1.1.0 1.3.6.1.4.1.43353.1.1.2.0 1.3.6.1.4.1.43353.1.1.3.0
[get][bmon@192.168.12.13][0] .1.3.6.1.4.1.43353.1.1.1.0 = Integer: 0
[get][bmon@192.168.12.13][1] .1.3.6.1.4.1.43353.1.1.2.0 = OctetString: SVSJm#2018/05/17-2023/05/16;SVSQm#2018/05/17-2023/05/16;156b1ba46762d0be#2018/05/17-2023/05/16
[get][bmon@192.168.12.13][2] .1.3.6.1.4.1.43353.1.1.3.0 = OctetString: 2021/1/14 14:0:48
[walk][bmon@192.168.12.13][0] .1.3.6.1.4.1.43353.1.1.1.0 = Integer: 0
[walk][bmon@192.168.12.13][0] .1.3.6.1.4.1.43353.1.1.2.0 = OctetString: SVSJm#2018/05/17-2023/05/16;SVSQm#2018/05/17-2023/05/16;156b1ba46762d0be#2018/05/17-2023/05/16
[walk][bmon@192.168.12.13][0] .1.3.6.1.4.1.43353.1.1.3.0 = OctetString: 2021/1/14 14:0:48
```

start snmp trap server:

```sh
$ snmp -trap :9162                                                                                     
2021/01/14 14:01:32 got trapdata from 127.0.0.1
[trap][127.0.0.1:65357][0] .1.3.6.1.2.1.1.3.0 = TimeTicks: 88396648
[trap][127.0.0.1:65357][1] .1.3.6.1.6.3.1.1.4.1.0 = ObjectIdentifier: .1.3.6.1.4.1.8072.2.3.0.1
[trap][127.0.0.1:65357][2] .1.3.6.1.4.1.43353.1.1.2.0 = OctetString: bingoohuang
2021/01/14 14:02:00 got trapdata from 127.0.0.1
[trap][127.0.0.1:53713][0] .1.3.6.1.2.1.1.3.0 = TimeTicks: 88399437
[trap][127.0.0.1:53713][1] .1.3.6.1.6.3.1.1.4.1.0 = ObjectIdentifier: .1.3.6.1.4.1.8072.2.3.0.1
[trap][127.0.0.1:53713][2] .1.3.6.1.4.1.8072.2.3.2.1 = Integer: 123456

$ snmp -trap :9162
2021/01/14 13:53:49 got trapdata from 127.0.0.1
[trap][127.0.0.1:59549][0] .1.3.6.1.2.1.1.3.0 = TimeTicks: 1610603629
[trap][127.0.0.1:59549][1] .1.3.6.1.2.1.1.6 = ObjectIdentifier: .1.3.6.1.2.1.1.6.10
[trap][127.0.0.1:59549][2] .1.3.6.1.2.1.1.7 = OctetString: Testing TCP trap...
[trap][127.0.0.1:59549][3] .1.3.6.1.2.1.1.8 = Integer: 123
```

Send A Test Trap:

```sh
$ snmptrap -v 2c -c public localhost:9162 '' 1.3.6.1.4.1.8072.2.3.0.1 1.3.6.1.4.1.43353.1.1.2.0  s bingoohuang
$ snmptrap -v 2c -c public localhost:9162 '' 1.3.6.1.4.1.8072.2.3.0.1 1.3.6.1.4.1.8072.2.3.2.1 i 123456
$ snmp -mode trapsend -t 127.0.0.1:9162
$ snmptranslate .1.3.6.1.2.1.1.3.0
DISMAN-EVENT-MIB::sysUpTimeInstance
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
