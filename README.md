# snmp

snmp package and tools in golang.

```sh
$ snmp -h
Usage: snmp [options] oids...
  -m     method to operate, get/walk/trapsend/translate (default is get/walk)
  -c     string Default SNMP community (default "public")
  -t     one or more SNMP targets (eg. -t 192.168.1.1 -t myCommunity@192.168.1.2:1234)
  -x/y/z one or more x/y/z vars (eg. -x 1-3 -y 1,3,5 -z 1,2-5)
  -o     one or more oids
  -s     trap server listening address(eg. :9162)
  -V     Verbose logging of packets, oid units, oid description and etc.
```

## examples

uses the SNMP GET/BULKWALK request to query for information on a network entity

```sh
❯ snmp -t bj@192.168.1.1 1.3.6.1.4.1.43353.1.1.1.0 1.3.6.1.4.1.43353.1.1.2.0 1.3.6.1.4.1.43353.1.1.3.0                      
[get][0][BJSER-MIB::MasterProcessStatus.0][.1.3.6.1.4.1.43353.1.1.1.0] = Integer: 0
[get][1][BJSER-MIB::CertificateValidate.0][.1.3.6.1.4.1.43353.1.1.2.0] = OctetString: SVSJm#2018/05/17-2023/05/16;SVSQm#2018/05/17-2023/05/16;156b1ba46762d0be#2018/05/17-2023/05/16
[get][2][BJSER-MIB::ServerTime.0][.1.3.6.1.4.1.43353.1.1.3.0] = OctetString: 2021/1/15 17:43:13
[walk][0][BJSER-MIB::MasterProcessStatus.0][.1.3.6.1.4.1.43353.1.1.1.0] = Integer: 0
[walk][0][BJSER-MIB::CertificateValidate.0][.1.3.6.1.4.1.43353.1.1.2.0] = OctetString: SVSJm#2018/05/17-2023/05/16;SVSQm#2018/05/17-2023/05/16;156b1ba46762d0be#2018/05/17-2023/05/16
[walk][0][BJSER-MIB::ServerTime.0][.1.3.6.1.4.1.43353.1.1.3.0] = OctetString: 2021/1/15 17:43:13
```

use `x` as a placeholder:

```sh
❯ snmp -m get -t bj@192.168.1.1 -o 1.3.6.1.4.1.43353.1.1.x.0 -x 1-3                                        
[0][BJSER-MIB::MasterProcessStatus.0][.1.3.6.1.4.1.43353.1.1.1.0] = Integer: 0
[1][BJSER-MIB::CertificateValidate.0][.1.3.6.1.4.1.43353.1.1.2.0] = OctetString: SVSJm#2018/05/17-2023/05/16;SVSQm#2018/05/17-2023/05/16;156b1ba46762d0be#2018/05/17-2023/05/16
[2][BJSER-MIB::ServerTime.0][.1.3.6.1.4.1.43353.1.1.3.0] = OctetString: 2021/1/15 17:44:26
```

start snmp trap server:

```sh
$ snmp -s :9162
2021/01/15 17:49:20 got trapdata from 127.0.0.1
[trap][0][DISMAN-EVENT-MIB::sysUpTimeInstance][.1.3.6.1.2.1.1.3.0] = TimeTicks: 98403325
[trap][1][SNMPv2-MIB::snmpTrapOID.0][.1.3.6.1.6.3.1.1.4.1.0] = ObjectIdentifier: .1.3.6.1.4.1.8072.2.3.0.1
[trap][2][BJSER-MIB::CertificateValidate.0][.1.3.6.1.4.1.43353.1.1.2.0] = OctetString: bingoohuang
2021/01/15 17:49:36 got trapdata from 127.0.0.1
[trap][0][DISMAN-EVENT-MIB::sysUpTimeInstance][.1.3.6.1.2.1.1.3.0] = TimeTicks: 98404969
[trap][1][SNMPv2-MIB::snmpTrapOID.0][.1.3.6.1.6.3.1.1.4.1.0] = ObjectIdentifier: .1.3.6.1.4.1.8072.2.3.0.1
[trap][2][NET-SNMP-EXAMPLES-MIB::netSnmpExampleHeartbeatRate][.1.3.6.1.4.1.8072.2.3.2.1] = Integer: 123456
2021/01/15 17:49:50 got trapdata from 127.0.0.1
[trap][0][DISMAN-EVENT-MIB::sysUpTimeInstance][.1.3.6.1.2.1.1.3.0] = TimeTicks: 1610704190
[trap][1][SNMPv2-MIB::sysLocation][.1.3.6.1.2.1.1.6] = ObjectIdentifier: .1.3.6.1.2.1.1.6.10
[trap][2][SNMPv2-MIB::sysServices][.1.3.6.1.2.1.1.7] = OctetString: Testing TCP trap...
[trap][3][SNMPv2-MIB::sysORLastChange][.1.3.6.1.2.1.1.8] = Integer: 123
```

Send A Test Trap:

```sh
$ snmptrap -v 2c -c public localhost:9162 '' 1.3.6.1.4.1.8072.2.3.0.1 1.3.6.1.4.1.43353.1.1.2.0  s bingoohuang
$ snmptrap -v 2c -c public localhost:9162 '' 1.3.6.1.4.1.8072.2.3.0.1 1.3.6.1.4.1.8072.2.3.2.1 i 123456
$ snmp -m trapsend -t 127.0.0.1:9162
```

Translate MIB OID names between numeric and textual forms like `snmptranslate`

```sh
$ snmp -m translate UCD-SNMP-MIB::dskAvail.1 .1.3.6.1.2.1.1.6                       
UCD-SNMP-MIB::dskAvail.1 => 1.3.6.1.4.1.2021.9.1.7.1
.1.3.6.1.2.1.1.6 => SNMPv2-MIB::sysLocation
```

```sh
$ snmp -V -m translate -o 1.3.6.1.4.1.2021.x -x 11.9.0,4.5.0,4.6.0,4.14.0,9.1.6.1,9.1.8.1,9.1.7.1 -o 1.3.6.1.4.1.43353.1.1.y.0 -y 1-3
2021/01/16 12:18:14 Oids:[1.3.6.1.4.1.2021.11.9.0 1.3.6.1.4.1.2021.4.5.0 1.3.6.1.4.1.2021.4.6.0 1.3.6.1.4.1.2021.4.14.0 1.3.6.1.4.1.2021.9.1.6.1 1.3.6.1.4.1.2021.9.1.8.1 1.3.6.1.4.1.2021.9.1.7.1 1.3.6.1.4.1.43353.1.1.1.0 1.3.6.1.4.1.43353.1.1.2.0 1.3.6.1.4.1.43353.1.1.3.0]
OidName: UCD-SNMP-MIB::ssCpuUser.0
Description: The percentage of CPU time spent processinguser-level code, calculated over the last minute.This object has been deprecated in favour of'ssCpuRawUser(50)', which can be used to calculatethe same metric, but over any desired time period.
1.3.6.1.4.1.2021.11.9.0 => UCD-SNMP-MIB::ssCpuUser.0
OidName: UCD-SNMP-MIB::memTotalReal.0 Unit: kB
Description: The total amount of real/physical memory installedon this host.
1.3.6.1.4.1.2021.4.5.0 => UCD-SNMP-MIB::memTotalReal.0
OidName: UCD-SNMP-MIB::memAvailReal.0 Unit: kB
Description: The amount of real/physical memory currently unusedor available.
1.3.6.1.4.1.2021.4.6.0 => UCD-SNMP-MIB::memAvailReal.0
OidName: UCD-SNMP-MIB::memBuffer.0 Unit: kB
Description: The total amount of real or virtual memory currentlyallocated for use as memory buffers.This object will not be implemented on hosts where theunderlying operating system does not explicitly identifymemory as specifically reserved for this purpose.
1.3.6.1.4.1.2021.4.14.0 => UCD-SNMP-MIB::memBuffer.0
OidName: UCD-SNMP-MIB::dskTotal.1
Description: Total size of the disk/partion (kBytes).For large disks (>2Tb), this value willlatch at INT32_MAX (2147483647).
1.3.6.1.4.1.2021.9.1.6.1 => UCD-SNMP-MIB::dskTotal.1
OidName: UCD-SNMP-MIB::dskUsed.1
Description: Used space on the disk.For large heavily-used disks (>2Tb), thisvalue will latch at INT32_MAX (2147483647).
1.3.6.1.4.1.2021.9.1.8.1 => UCD-SNMP-MIB::dskUsed.1
OidName: UCD-SNMP-MIB::dskAvail.1
Description: Available space on the disk.For large lightly-used disks (>2Tb), thisvalue will latch at INT32_MAX (2147483647).
1.3.6.1.4.1.2021.9.1.7.1 => UCD-SNMP-MIB::dskAvail.1
1.3.6.1.4.1.43353.1.1.1.0 => BJSER-MIB-MIB::MasterProcessStatus.0
1.3.6.1.4.1.43353.1.1.2.0 => BJSER-MIB-MIB::CertificateValidate.0
1.3.6.1.4.1.43353.1.1.3.0 => BJSER-MIB-MIB::ServerTime.0
```

## resources

### public resources

1. [SMPP (3.4) Client Library in pure Go](https://github.com/linxGnu/gosmpp), This library is well tested with SMSC simulators: [Melroselabs SMSC](https://melroselabs.com/services/smsc-simulator/#smsc-simulator-try)

### SNMP v2 Trap

[SNMP Trap - How To Send A Test Trap](https://support.nagios.com/kb/article.php?id=493)

* Command form: `snmptrap -v <snmp_version> -c <community> <destination_host> <uptime> <OID_or_MIB> <object> <value_type> <value>`
* Using MIB: `snmptrap -v2c -c public localhost '' NET-SNMP-EXAMPLES-MIB::netSnmpExampleHeartbeatNotification netSnmpExampleHeartbeatRate i 123456`
* Shortening MIB: `snmptrap -v2c -c public localhost '' netSnmpExampleHeartbeatNotification netSnmpExampleHeartbeatRate i 123456`
* Using OID: `snmptrap -v 2c -c public localhost '' 1.3.6.1.4.1.8072.2.3.0.1 1.3.6.1.4.1.8072.2.3.2.1 i 123456`

The commands above required the following settings in /etc/snmp/snmptrapd.conf

  disableAuthorization yes
  traphandle default /usr/sbin/snmptthandler

### SNMP定义名词术语

[SNMP定义](https://github.com/fenggolang/collect)

* SNMP：Simple Network Management Protocol(简单网络管理协议)，是一个标准的用于管理基于IP网络上设备的协议。
  * SNMP的主要功能: 通过应答POLLING(轮询)来反馈当前设备状态;
  * SNMP的工作方式: 管理员需要向设备获取数据,所以SNMP提供了"读"操作;管理员需要向设备执行设置操作,所以SNMP提供了"写"操作; 设备需要在重要状况改变的时候,向管理员通报事件的发生,所以SNMP提供了"Trap" 操作;
  * SNMP被设计为工作在TCP/IP协议族上.SNMP基于TCP/IP协议工作,对网络中支持SNMP协议的设备进行管理.所有支持SNMP协议的设备都提供SNMP这个统一界面，使得管理员可以使用统一的操作进行管理，而不必理会设备是什么类型、是哪个厂家生产的.
* OID: 对象标识符（OID－Object Identifiers），是SNMP代理提供的具有唯一标识的键值，MIB（管理信息基）提供数字化OID到可读文本的映射。
  * SNMP OID是用一种按照层次化格式组织的、树状结构中的唯一地址来表示的，它与DNS层次相似。 
  * sample structure of an OID: 
    * `iso(1).org(3).dod(6).internet(1).private(4).transition(868).products(2).chassis(4).card(1).slotCps(2)-cpsSlotSummary(1).cpsModuleTable(1).cpsModuleEntry(1).cpsModuleModel(3).3562.3`
    * or just: `1.3.6.1.4.868.2.4.1.2.1.1.1.3.3562.3`
  * ![image](https://user-images.githubusercontent.com/1940588/104560584-0a639380-5681-11eb-8de8-a6f71b8788c9.png)
  * [oidref.com](https://oidref.com/1.3.6.1.6.3.1.1.4.1)
  * [http://oid-info.com/](http://oid-info.com/get/1.3.6.1.4.1.2021.4.5)
* MIB：Management Information Base(管理信息库)，定义代理进程中所有可被查询和修改的参数。
  * 一个OID看起来和一个IPv6的地址很象，并且不同的厂商有不同的前缀等信息。
  * OID都非常长，使得人们难以记住，或者对他非常感冒。因此，人们就设计了一种将数字OID翻译为人们可读的格式。这种翻译映射被保存在一个被称为 “管理信息基础"MIB的、可传递的无格式文本文件里。
  * 使用SNMP或者向SNMP设备查询，你不需要使用MIB，但是，如果没有MIB，你就得猜测你正在查看的数据是什么。
  * 某些情况下，不使用MIB也非常简单，例如查看主机名、磁盘使用率数字，或者端口状态信息。其他情况下，就非常困难了，这个时候使用MIB就非常有帮助。
  * 对于准备编写的应用程序来说，为了让用户避免妥当安装MIB带来的麻烦，而严格使用数字OID很常见。
  * 安装一个MIB的动作，只是将他放置到你的SNMP客户端应用软件能够搜索到并进行上述翻译映射工作的某个位置而已(/usr/share/snmp/mibs, ~/.snmp/mibs)。
* SMI：Structure of Management Information(管理信息结构)
  * SMI定义了SNMP中使用到的ASN.1类型、语法，并定义了SNMP中使用到的类型、宏、符号等。SMI用于后续协议的描述和MIB的定义。每个版本的SNMP都可能定义自己的SMI。
  * [python parse MIB files from ASN.1 SMI sources](https://github.com/qmsk/snmpbot/tree/master/scripts)
  * [MIB json example](https://github.com/qmsk/snmpbot/blob/master/mibs/test/TEST2-MIB.json)
* SNMP的常用版本有三个：SNMPv1、SNMPv2、SNMPv3
  * SNMPv1是为基于公共管理的初始标准(1988年）。
  * SNMPv2是SNMPv1框架下衍生出来的，但是没有定义信息，其后修订为SNMPv2c（1908年）,一个带有于SNMPv1类似信息格式的给予公共管理的版本。
  * SNMPv2添加了几个新的数据类型（Counter32、Counter64、Gauge32、UInteger32、NsapAdress 以及BIT STRING),以及对OID表和OID值的设置的增强。
  * SNMPv3（2004年）是一个带有新的信息格式、ACL、安全功能和远处SNMP参数配置的、扩展了SNMPv2框架的版本。
  * SNMPv1和SNMPv2被广泛应用，但是由于这些协议的不安全特性，通常只使用只读访问。通常，除非你确实需要附加安全特性，否则你不需要过多地关注SNMPv3。SNMPv3是具有安全性的通信协议。
* SNMP是基于几个其他规范的，包括定义给予SNMP的数据类型的ASN.1 BER(Abstract Syntax Notation 1 Basic Encoding Rules), 
* 以及详细描述有SNMP MIB使用的语法的管理信息结构（SMI）。**如果你选择修改MIB的某些东西，你需要学习SMIv2和ASN.1语法**，不过其他情况下你只需要在兴趣时看看他，而不必学习他。

### 安装使用 SNMP

[安装使用 SNMP](https://github.com/fenggolang/collect)

```sh
# 安装
yum install net-snmp net-snmp-utils net-snmp* -y

# 配置
vim /etc/snmp/snmpd.conf com2sec notConfigUser default ccssoft view all included .1 access notConfigGroup ""      any
noauth exact all none none includeAllDisks rocommunity ccssoft disk / disk /home

# 启动snmp服务
systemctl enable snmpd systemctl start snmpd

# 确保iptables防火墙开放了udp 161端口的访问权限
# 配置防火墙规则运行snmp端口161
iptables -A INPUT -p tcp -m state --state NEW -m tcp --dport 161 -j ACCEPT iptables -A INPUT -p udp -m state --state
NEW -m udp --dport 161 -j ACCEPT systemctl start snmpd systemctl enable snmpd iptables -D INPUT -p tcp -m state
--state NEW -m tcp --dport 161 -j ACCEPT iptables -D INPUT -p udp -m state --state NEW -m udp --dport 161 -j ACCEPT
iptables -I INPUT -p tcp -m state --state NEW -m tcp --dport 161 -j ACCEPT iptables -I INPUT -p udp -m state --state
NEW -m udp --dport 161 -j ACCEPT

# 如果是firewalld则如下方式添加
firewall-cmd --zone=public --add-port=161/tcp --permanent firewall-cmd --zone=public --add-port=161/udp --permanent
firewall-cmd --reload systemctl restart snmpd

# 查看所有开放的端口
firewall-cmd --zone=public --list-ports

# 查看snmp版本
[root@paas ~]# snmpget --version NET-SNMP version: 5.7.2
[root@paas ~]#

# 查看一下安装的snmp软件包
rpm -qa | grep net-snmp*

snmpget -c ccssoft -v 2c localhost .1.3.6.1.4.1.2021.11.9.0

# snmpget 模拟snmp的GetRequest操作的工具。用来获取一个或几个管理信息。用来读取管理信息的内容。
# 获取设备的描述信息
[root@paas ~]# snmpget -c ccssoft -v 2c paas-node1.m8.ccs sysDescr.0 SNMPv2-MIB::sysDescr.0 = STRING: Linux
paas-node1.m8.ccs 3.10.0-514.26.2.el7.x86_64 #1 SMP Tue Jul 4 15:04:05 UTC 2017 x86_64
[root@paas ~]# uname -a Linux paas.m8.ccs 3.10.0-693.5.2.el7.x86_64 #1 SMP Fri Oct 20 20:32:50 UTC 2017 x86_64 x86_64
x86_64 GNU/Linux
[root@paas ~]#

# 获取磁盘信息
[root@paas ~]# snmpdf -v2c -c ccssoft localhost
```

### snmpwalk和snmpget的区别

[snmpwalk和snmpget的区别](https://github.com/fenggolang/collect)

snmpwalk是对OID值的遍历（比如某个OID值下面有N个节点，则依次遍历出这N个节点的值。如果对某个叶子节点的OID值做walk，则取得到数据就不正确了，因为它会认为该节点是某些节点的父节点，
而对其进行遍历，而实际上该节点已经没有子节点了，那么它会取出与该叶子节点平级的下一个叶子节点的值，而不是当前请求的节子节点的值。）

snmpget是取具体的OID的值。（适用于OID值是一个叶子节点的情况）

### SNMP监控一些常用OID的总结

处理器、存储器、存储设备、网络环境为最常被监控的项目，[这些项目的细节信息如下表](https://zh.wikipedia.org/wiki/%E7%AE%80%E5%8D%95%E7%BD%91%E7%BB%9C%E7%AE%A1%E7%90%86%E5%8D%8F%E8%AE%AE)

|     项目      | 访问方式 |            OID            |                 MIB                 |
|-------------|------|---------------------------|-------------------------------------|
|     处理器     | 处理器  |            处理器            |                 处理器                 |
|   用户占用时间比   | GET  | .1.3.6.1.4.1.2021.11.9.0  |       UCD-SNMP-MIB::ssCpuUser       |
|   系统占用时间比   | GET  | .1.3.6.1.4.1.2021.11.10.0 |      UCD-SNMP-MIB::ssCpuSystem      |
|    闲置时间比    | GET  | .1.3.6.1.4.1.2021.11.11.0 |       UCD-SNMP-MIB::ssCpuIdle       |
| 每个 Core 的用量 | WALK |  .1.3.6.1.2.1.25.3.3.1.2  | HOST-RESOURCES-MIB::hrProcessorLoad |
|     存储器     | 存储器  |            存储器            |                 存储器                 |
|    存储器容量    | GET  |  .1.3.6.1.4.1.2021.4.5.0  |     UCD-SNMP-MIB::memTotalReal      |
|   存储器消耗量    | GET  |  .1.3.6.1.4.1.2021.4.6.0  |     UCD-SNMP-MIB::memAvailReal      |
|   存储器剩余量    | GET  | .1.3.6.1.4.1.2021.4.11.0  |     UCD-SNMP-MIB::memTotalFree      |
|   虚拟内存容量    | GET  |  .1.3.6.1.4.1.2021.4.3.0  |     UCD-SNMP-MIB::memTotalSwap      |
|   虚拟内存剩余量   | GET  |  .1.3.6.1.4.1.2021.4.4.0  |     UCD-SNMP-MIB::memAvailSwap      |
|    存储设备     | 存储设备 |           存储设备            |                存储设备                 |
|    各磁盘容量    | WALK | .1.3.6.1.4.1.2021.9.1.6.1 |       UCD-SNMP-MIB::dskTotal        |
|   各磁盘消耗量    | WALK | .1.3.6.1.4.1.2021.9.1.7.1 |       UCD-SNMP-MIB::dskAvail        |
|  各磁盘消耗量百分比  | WALK | .1.3.6.1.4.1.2021.9.1.9.1 |      UCD-SNMP-MIB::dskPercent       |
|    网络环境     | 网络环境 |           网络环境            |                网络环境                 |
|   网络设备名称    | WALK |  .1.3.6.1.2.1.31.1.1.1.1  |           IF-MIB::ifName            |


[from](https://www.cnblogs.com/aspx-net/p/3554044.html)

| 系统参数  OID               | 描述        | 备注                | 请求方式 |
|-------------------------|-----------|-------------------|------|
| .1.3.6.1.2.1.1.1.0      | 获取系统基本信息  | SysDesc           | GET  |
| .1.3.6.1.2.1.1.3.0      | 监控时间      | sysUptime         | GET  |
| .1.3.6.1.2.1.1.4.0      | 系统联系人     | sysContact        | GET  |
| .1.3.6.1.2.1.1.5.0      | 获取机器名     | SysName           | GET  |
| .1.3.6.1.2.1.1.6.0      | 机器坐在位置    | SysLocation       | GET  |
| .1.3.6.1.2.1.1.7.0      | 机器提供的服务   | SysService        | GET  |
| .1.3.6.1.2.1.25.4.2.1.2 | 系统运行的进程列表 | hrSWRunName       | WALK |
| .1.3.6.1.2.1.25.6.3.1.2 | 系统安装的软件列表 | hrSWInstalledName | WALK |

| 网络接口  OID             | 描述                 | 备注             | 请求方式         |
|-----------------------|--------------------|----------------|--------------|
| .1.3.6.1.2.1.2.1.0    | 网络接口的数目            | IfNumber       | GET          |
| .1.3.6.1.2.1.2.2.1.2  | 网络接口信息描述           | IfDescr        | WALK         |
| .1.3.6.1.2.1.2.2.1.3  | 网络接口类型             | IfType         | WALK         |
| .1.3.6.1.2.1.2.2.1.4  | 接口收发的最大IP数据报[BYTE] | IfMTU          | WALK         |
| .1.3.6.1.2.1.2.2.1.5  | 接口当前带宽[bps]        | IfSpeed        | WALK         |
| .1.3.6.1.2.1.2.2.1.6  | 接口的物理地址            | IfPhysAddress  | WALK         |
| .1.3.6.1.2.1.2.2.1.8  | 接口当前操作状态[up        | down]          | IfOperStatus |        WALK         |
| .1.3.6.1.2.1.2.2.1.10 | 接口收到的字节数           | IfInOctet      | WALK         |
| .1.3.6.1.2.1.2.2.1.16 | 接口发送的字节数           | IfOutOctet     | WALK         |
| .1.3.6.1.2.1.2.2.1.11 | 接口收到的数据包个数         | IfInUcastPkts  | WALK         |
| .1.3.6.1.2.1.2.2.1.17 | 接口发送的数据包个数         | IfOutUcastPkts | WALK         |

| CPU及负载 OID                 | 描述                 | 备注                | 请求方式 |
|----------------------------|--------------------|-------------------|------|
| .1.3.6.1.4.1.2021.11.9.0   | 用户CPU百分比           | ssCpuUser         | GET  |
| .1.3.6.1.4.1.2021.11.10.0  | 系统CPU百分比           | ssCpuSystem       | GET  |
| .1.3.6.1.4.1.2021.11.11.0  | 空闲CPU百分比           | ssCpuIdle         | GET  |
| .1.3.6.1.4.1.2021.11.50.0  | 原始用户CPU使用时间        | ssCpuRawUser      | GET  |
| .1.3.6.1.4.1.2021.11.51.0  | 原始nice占用时间         | ssCpuRawNice      | GET  |
| .1.3.6.1.4.1.2021.11.52.0  | 原始系统CPU使用时间        | ssCpuRawSystem.   | GET  |
| .1.3.6.1.4.1.2021.11.53.0  | 原始CPU空闲时间          | ssCpuRawIdle      | GET  |
| .1.3.6.1.2.1.25.3.3.1.2    | CPU的当前负载，N个核就有N个负载 | hrProcessorLoad   | WALK |
| .1.3.6.1.4.1.2021.11.3.0   | -                  | ssSwapIn          | GET  |
| .1.3.6.1.4.1.2021.11.4.0   | -                  | SsSwapOut         | GET  |
| .1.3.6.1.4.1.2021.11.5.0   | -                  | ssIOSent          | GET  |
| .1.3.6.1.4.1.2021.11.6.0   | -                  | ssIOReceive       | GET  |
| .1.3.6.1.4.1.2021.11.7.0   | -                  | ssSysInterrupts   | GET  |
| .1.3.6.1.4.1.2021.11.8.0   | -                  | ssSysContext      | GET  |
| .1.3.6.1.4.1.2021.11.54.0  | -                  | ssCpuRawWait      | GET  |
| .1.3.6.1.4.1.2021.11.56.0  | -                  | ssCpuRawInterrupt | GET  |
| .1.3.6.1.4.1.2021.11.57.0  | -                  | ssIORawSent       | GET  |
| .1.3.6.1.4.1.2021.11.58.0  | -                  | ssIORawReceived   | GET  |
| .1.3.6.1.4.1.2021.11.59.0  | -                  | ssRawInterrupts   | GET  |
| .1.3.6.1.4.1.2021.11.60.0  | -                  | ssRawContexts     | GET  |
| .1.3.6.1.4.1.2021.11.61.0  | -                  | ssCpuRawSoftIRQ   | GET  |
| .1.3.6.1.4.1.2021.11.62.0  | -                  | ssRawSwapIn.      | GET  |
| .1.3.6.1.4.1.2021.11.63.0  | -                  | ssRawSwapOut      | GET  |
| .1.3.6.1.4.1.2021.10.1.3.1 | -                  | Load5             | GET  |
| .1.3.6.1.4.1.2021.10.1.3.2 | -                  | Load10            | GET  |
| .1.3.6.1.4.1.2021.10.1.3.3 | -                  | Load15            | GET  |

| 内存及磁盘    OID             | 描述                                      | 备注                       | 请求方式 |
|--------------------------|-----------------------------------------|--------------------------|------|
| .1.3.6.1.2.1.25.2.2.0    | 获取内存大小                                  | hrMemorySize             | GET  |
| .1.3.6.1.2.1.25.2.3.1.1  | 存储设备编号                                  | hrStorageIndex           | WALK |
| .1.3.6.1.2.1.25.2.3.1.2  | 存储设备类型                                  | hrStorageType[OID]       | WALK |
| .1.3.6.1.2.1.25.2.3.1.3  | 存储设备描述                                  | hrStorageDescr           | WALK |
| .1.3.6.1.2.1.25.2.3.1.4  | 簇的大小                                    | hrStorageAllocationUnits | WALK |
| .1.3.6.1.2.1.25.2.3.1.5  | 簇的的数目                                   | hrStorageSize            | WALK |
| .1.3.6.1.2.1.25.2.3.1.6  | 使用多少，跟总容量相除就是占用率                        | hrStorageUsed            | WALK |
| .1.3.6.1.4.1.2021.4.3.0  | Total Swap Size(虚拟内存)                   | memTotalSwap             | GET  |
| .1.3.6.1.4.1.2021.4.4.0  | Available Swap Space                    | memAvailSwap             | GET  |
| .1.3.6.1.4.1.2021.4.5.0  | Total RAM in machine                    | memTotalReal             | GET  |
| .1.3.6.1.4.1.2021.4.6.0  | Total RAM used                          | memAvailReal             | GET  |
| .1.3.6.1.4.1.2021.4.11.0 | Total RAM Free                          | memTotalFree             | GET  |
| .1.3.6.1.4.1.2021.4.13.0 | Total RAM Shared                        | memShared                | GET  |
| .1.3.6.1.4.1.2021.4.14.0 | Total RAM Buffered                      | memBuffer                | GET  |
| .1.3.6.1.4.1.2021.4.15.0 | Total Cached Memory                     | memCached                | GET  |
| .1.3.6.1.4.1.2021.9.1.2  | Path where the disk is mounted          | dskPath                  | WALK |
| .1.3.6.1.4.1.2021.9.1.3  | Path of the device for the partition    | dskDevice                | WALK |
| .1.3.6.1.4.1.2021.9.1.6  | Total size of the disk/partion (kBytes) | dskTotal                 | WALK |
| .1.3.6.1.4.1.2021.9.1.7  | Available space on the disk             | dskAvail                 | WALK |
| .1.3.6.1.4.1.2021.9.1.8  | Used space on the disk                  | dskUsed                  | WALK |
| .1.3.6.1.4.1.2021.9.1.9  | Percentage of space used on disk        | dskPercent               | WALK |
| .1.3.6.1.4.1.2021.9.1.10 | Percentage of inodes used on disk       | dskPercentNode           | WALK |

### OID要不要以点开头

from [here](https://support.microfocus.com/kb/doc.php?id=7743528)

There is a distinction between those specified with a leading dot (i.e. '.1.3.6.1.2.1.1.3.0') 
and those without (i.e. '1.3.0').

> If an OID has a leading dot, it is assumed the OID is fully qualified. 

> If there is no leading dot, it is assumed that the OID is prefixed with 'iso.org.dod.internet.mgmt.mib'.
 
In the examples above both '.1.3.6.1.2.1.1.3.0' and '1.3.0' are equivalent to 'sysUpTime.0'.

### 一些图

来自[这里](http://gosnmpapi.webnms.com/snmpget)

![image](https://user-images.githubusercontent.com/1940588/104563024-1bfa6a80-5684-11eb-8652-319b836b41d5.png)

![image](https://user-images.githubusercontent.com/1940588/104563388-8ca18700-5684-11eb-949b-7e34b5eae7b9.png)

![image](https://user-images.githubusercontent.com/1940588/104563414-94612b80-5684-11eb-8569-7598e53acaac.png)

[Monitoring via SNMP](https://www.paessler.com/info/snmp_mibs_and_oids_an_overview)

![image](https://user-images.githubusercontent.com/1940588/104808403-f8b8f200-5820-11eb-8d05-1420d1e9e109.png)


### Configure SNMP service on Mac OSX

1. `sudo -i`
2. `vi /etc/snmp/snmpd.conf`
3. replace
  
  ```
  com2sec local localhost COMMUNITY
  com2sec mynetwork NETWORK/24 COMMUNITY
  ```

  with
  ```
  com2sec local localhost private
  com2sec mynetwork NETWORK/24 public
  ```

1. replace `rocommunity public default .1.3.6.1.2.1.1.4` with `rocommunity public default .1`
2. uncomment `#rwcommunity private`
3. `launchctl unload /System/Library/LaunchDaemons/org.net-snmp.snmpd.plist`
4. `launchctl load -w /System/Library/LaunchDaemons/org.net-snmp.snmpd.plist`
5. test

  ```sh
  $ snmp -m get -t 127.0.0.1 -o 1.3.6.1.4.1.2021.x -x 11.9.0,4.5.0,4.6.0,4.14.0,9.1.6.1,9.1.8.1,9.1.7.1
  [0][UCD-SNMP-MIB::ssCpuUser.0][.1.3.6.1.4.1.2021.11.9.0] = Integer: 3
  [1][UCD-SNMP-MIB::memTotalReal.0][.1.3.6.1.4.1.2021.4.5.0] = Integer: 16777216
  [2][UCD-SNMP-MIB::memAvailReal.0][.1.3.6.1.4.1.2021.4.6.0] = Integer: 2596108
  [3][UCD-SNMP-MIB::memBuffer.0][.1.3.6.1.4.1.2021.4.14.0] = NoSuchObject: <nil>
  [4][UCD-SNMP-MIB::dskTotal.1][.1.3.6.1.4.1.2021.9.1.6.1] = Integer: 488245280
  [5][UCD-SNMP-MIB::dskUsed.1][.1.3.6.1.4.1.2021.9.1.8.1] = Integer: 14694344
  [6][UCD-SNMP-MIB::dskAvail.1][.1.3.6.1.4.1.2021.9.1.7.1] = Integer: 212321536
  ```

### Mib Browser

MIB浏览器是工程师管理启用SNMP的网络设备和应用程序不可或缺的工具。 它允许用户加载标准的，专有的MIB，甚至一些错误的MIB。 它还允许它们发出SNMP请求以检索代理的数据，或对代理进行更改。 内置陷阱接收器可以根据其规则引擎接收和处理SNMP陷阱。

#### iReasoning MIB浏览器

[ireasoning Mib Browser](http://www.ireasoning.com/download.shtml)

一个强大和易于使用的工具由iReasoning SNMP API提供支持。

![image](https://user-images.githubusercontent.com/1940588/104807784-82b28c00-581c-11eb-91cc-f56531c21ebb.png)

#### IPNetMonitorX

IPNetMonitorX将直观的Mac用户界面与深入的OS X知识结合在一起，从一个简单的应用程序中提供重要的问题解决细节。如果你曾经为什么一些互联网功能不能正常工作，或者只是想更好地探索和理解，IPNetMonitorX在这里提供帮助。每个工具的内置描述都会澄清底层网络概念，因为它会指导您回答有关Internet如何工作的实际问题。

https://www.osxwin.com/s/ipnetmonitorx
http://www.sustworks.com/rw_site/download.php

主要特性：
一个完整的Internet工具包。一致的集成设计提供了速度和效率。
浮动工具调色板使您能够快速查看和访问所需的工具。每个工具的布局都针对手头的任务进行了优化。
全面的显示显示了TCP / IP的实际应用，非常适合调整网络性能。
快速，异步，多线程的设计使您能够在发生网络行为时看到。
从单个应用程序中同时打开并使用同一工具的多个工具或实例。
输入字段支持最近的目标或历史记录，自动完成和拖放。自动将数据从一个工具传输到另一个工具。
不需要以root身份登录。
工具：
地址扫描 AirPort 信号 Bonjour扫描 连接列表DHCP租约 DHCP测试 DNS查找 DNS查询 DNS RBL 手指 界面信息链接率
监控 ping 端口扫描 服务器扫描 服务信息 SNMPd配置 SNMP工具 子网计算器 TCP信息 TCP转储 跟踪路由 Whois

![image](https://user-images.githubusercontent.com/1940588/104808016-39fbd280-581e-11eb-9132-4751e402c9a6.png)

![image](https://user-images.githubusercontent.com/1940588/104808043-90691100-581e-11eb-8b03-e1dc01245c83.png)

## logs

1. [gosnmpd](https://github.com/bingoohuang/gosnmpd)

```log
$ gosnmpd run
```

```log
$ snmp -m get -t 127.0.0.1:1161 -version 3 -userName testuser -authPassword testauth -authProtocol MD5 -contextName public -o 1.3.6.1.4.1.2021.9.1.6.1
[0][UCD-SNMP-MIB::dskTotal.1][.1.3.6.1.4.1.2021.9.1.6.1] => Integer: 80569
```

```log
$ snmp -m walk -t 127.0.0.1:1161 -version 3 -userName testuser -authPassword testauth -authProtocol MD5 -contextName public -o 1.3.6.1.4.1.2021
[0][UCD-SNMP-MIB::memIndex][.1.3.6.1.4.1.2021.4.1] => Integer: 1
[1][UCD-SNMP-MIB::memErrorName][.1.3.6.1.4.1.2021.4.2] => OctetString: swap
[2][UCD-SNMP-MIB::memTotalSwap][.1.3.6.1.4.1.2021.4.3] => Integer: 0
[3][UCD-SNMP-MIB::memAvailSwap][.1.3.6.1.4.1.2021.4.4] => Integer: 0
[4][UCD-SNMP-MIB::memTotalReal][.1.3.6.1.4.1.2021.4.5] => Integer: 3825904
[5][UCD-SNMP-MIB::memAvailReal][.1.3.6.1.4.1.2021.4.6] => Integer: 3122672
[6][UCD-SNMP-MIB::memTotalFree][.1.3.6.1.4.1.2021.4.11] => Integer: 3122672
[7][UCD-SNMP-MIB::memMinimumSwap][.1.3.6.1.4.1.2021.4.12] => Integer: 0
[8][UCD-SNMP-MIB::memBuffer][.1.3.6.1.4.1.2021.4.14] => Integer: 120372
[9][UCD-SNMP-MIB::memCached][.1.3.6.1.4.1.2021.4.15] => Integer: 2906880
[10][UCD-SNMP-MIB::memSwapError][.1.3.6.1.4.1.2021.4.100] => Integer: 0
[11][UCD-SNMP-MIB::memSwapErrorMsg][.1.3.6.1.4.1.2021.4.101] => OctetString:
[12][UCD-SNMP-MIB::dskIndex.1][.1.3.6.1.4.1.2021.9.1.1.1] => Integer: 1
[13][UCD-SNMP-MIB::dskPath.1][.1.3.6.1.4.1.2021.9.1.2.1] => OctetString: /
[14][UCD-SNMP-MIB::dskDevice.1][.1.3.6.1.4.1.2021.9.1.3.1] => OctetString: /
[15][UCD-SNMP-MIB::dskTotal.1][.1.3.6.1.4.1.2021.9.1.6.1] => Integer: 80569
[16][UCD-SNMP-MIB::dskAvail.1][.1.3.6.1.4.1.2021.9.1.7.1] => Integer: 39829
[17][UCD-SNMP-MIB::dskUsed.1][.1.3.6.1.4.1.2021.9.1.8.1] => Integer: 37365
[18][UCD-SNMP-MIB::dskPercent.1][.1.3.6.1.4.1.2021.9.1.9.1] => Integer: 48
[19][UCD-SNMP-MIB::laIndex.1][.1.3.6.1.4.1.2021.10.1.1.1] => Integer: 1
[20][UCD-SNMP-MIB::laIndex.2][.1.3.6.1.4.1.2021.10.1.1.2] => Integer: 2
[21][UCD-SNMP-MIB::laIndex.3][.1.3.6.1.4.1.2021.10.1.1.3] => Integer: 3
[22][UCD-SNMP-MIB::laNames.1][.1.3.6.1.4.1.2021.10.1.2.1] => OctetString: Load-1
[23][UCD-SNMP-MIB::laNames.2][.1.3.6.1.4.1.2021.10.1.2.2] => OctetString: Load-5
[24][UCD-SNMP-MIB::laNames.3][.1.3.6.1.4.1.2021.10.1.2.3] => OctetString: Load-15
[25][UCD-SNMP-MIB::laLoad.1][.1.3.6.1.4.1.2021.10.1.3.1] => OctetString: 0.03
[26][UCD-SNMP-MIB::laLoad.2][.1.3.6.1.4.1.2021.10.1.3.2] => OctetString: 0.03
[27][UCD-SNMP-MIB::laLoad.3][.1.3.6.1.4.1.2021.10.1.3.3] => OctetString: 0.01
[28][UCD-SNMP-MIB::laLoadInt.1][.1.3.6.1.4.1.2021.10.1.5.1] => Integer: 0
[29][UCD-SNMP-MIB::laLoadInt.2][.1.3.6.1.4.1.2021.10.1.5.2] => Integer: 0
[30][UCD-SNMP-MIB::laLoadInt.3][.1.3.6.1.4.1.2021.10.1.5.3] => Integer: 0
[31][UCD-SNMP-MIB::ssIndex][.1.3.6.1.4.1.2021.11.1] => Integer: 0
[32][UCD-SNMP-MIB::ssErrorName][.1.3.6.1.4.1.2021.11.2] => OctetString: systemStats
[33][UCD-SNMP-MIB::ssCpuRawUser][.1.3.6.1.4.1.2021.11.50] => Counter32: 22918
[34][UCD-SNMP-MIB::ssCpuRawNice][.1.3.6.1.4.1.2021.11.51] => Counter32: 617
[35][UCD-SNMP-MIB::ssCpuRawSystem][.1.3.6.1.4.1.2021.11.52] => Counter32: 20240
[36][UCD-SNMP-MIB::ssCpuRawIdle][.1.3.6.1.4.1.2021.11.53] => Counter32: 17642486
[37][UCD-SNMP-MIB::ssCpuRawWait][.1.3.6.1.4.1.2021.11.54] => Counter32: 3169
[38][UCD-SNMP-MIB::ssCpuRawInterrupt][.1.3.6.1.4.1.2021.11.56] => Counter32: 0
[39][UCD-SNMP-MIB::ssIORawSent][.1.3.6.1.4.1.2021.11.57] => Counter32: 26512384
[40][UCD-SNMP-MIB::ssIORawReceived][.1.3.6.1.4.1.2021.11.58] => Counter32: 575607
[41][UCD-SNMP-MIB::ssRawInterrupts][.1.3.6.1.4.1.2021.11.59] => Counter32: 1080348069
[42][UCD-SNMP-MIB::ssRawContexts][.1.3.6.1.4.1.2021.11.60] => Counter32: 1547188702
[43][UCD-SNMP-MIB::ssCpuRawSoftIRQ][.1.3.6.1.4.1.2021.11.61] => Counter32: 537
[44][UCD-SNMP-MIB::ssCpuRawSteal][.1.3.6.1.4.1.2021.11.64] => Counter32: 0
[45][UCD-SNMP-MIB::ssCpuRawGuest][.1.3.6.1.4.1.2021.11.65] => Counter32: 0
```