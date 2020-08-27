grammar zbxcsv;

fragment DIGIT: [0-9];
NUM: DIGIT+;
DOT: '.';
WORD: ~[,\r\n".+&:]+;
// 主机名称,显示名称,所属群组id,接口IP,模板ID,禁用主机,使用IP,DNS,描述

file: row+;

row:
	hostname ',' vis_name ',' group_id ',' ifaces ',' templ_id ',' disable_host ',' use_ip ',' dns
		',' description '\r'? EOF? '\n'? # host;

port: NUM;
ip: NUM DOT NUM DOT NUM DOT NUM;
iface_type: 'snmp' | 'ipmi' | 'jmx' | 'agent';
iface: ip ':' port '+' iface_type;
ifaces: iface ('&' iface)* # interfac;

dns: NUM DOT NUM DOT NUM DOT NUM;
hostname: WORD;
vis_name: WORD;
group_id: NUM;
groups: group_id ('&' group_id)+;
templ_id: NUM;
disable_host: NUM;
use_ip: NUM;
description: WORD;
