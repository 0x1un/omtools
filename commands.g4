grammar commands;

NL: '\n';
SPACE: ' ';
ALL: '.*?';
WORD: ~[ \r\n\t];
// 通用命令
EXIT: 'exit' | 'bye';
HELP: 'help';
RECON: 'recon';
GO: 'go';

// zabbix
ZBX: 'zbx';
LIST: 'list';
QUERY: 'query';
CFG: 'cfg';
HOST: 'host';
GROUP: 'group';
EXPORT: 'export';
JSON_XML: ALL '.' 'json' | ALL '.' 'xml';

// AD
ADD: 'add';
SINGLE: 'single';
USER: 'user';
FROM: 'from';
DEL: 'del';
WITH: 'with';
INFO: 'info';
BY: 'by';
DIS: 'dis';
ENA: 'ena';
UNLOCK: 'unlock';
AD: 'ad';
WS: [ \n]+ -> skip;

input: (';' value)*;

value: '(' WORD ')';