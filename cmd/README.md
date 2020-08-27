# omcmd

omcmd是一个命令行下的维护小工具，目前主要包含如下的功能：

- Zabbix
- Windows Active Direcotry

目前的功能主要集中在shell模式当中

```
git clone https://github.com/0x1un/omtools.git
cd omtools; cd cmd
go build .
./omcmd shell
omtools »  
```

输入help
```
omtools »  help // 输出命令树帮助
commands:
    query 
    ├── host 
    ├────── by 
    ├── tpl 
    ├────── by 
    ├── graph 
    ├────── by 
    ├── info 
    ├────── * 
    ├────── all 
    ├── user 
    ├────── by 
    dis 
    ena 
    unlock 
    re 
    ├── con 
    ├────── ad 
    ├────── zbx 
    add 
    ├── single 
    ├────── user 
    ├── user 
    ├────── from 
    del 
    ├── user 
    ├────── with 
    go 
    ├── zbx 
    ├── ad 
    list 
    ├── host 
    login 
    bye 
    help 
```

# command help
```shell
###### 通用 ######
bye # 退出shell
exit # 退出shell
help # 打印命令树
re con [zbx|ad] # 重新连接zbx或者ad
go [zbx|ad] # 进入zbx模式或者ad模式

###### AD ######
add single user # 添加单个用户
add user from [example.csv] # 从example.csv中获取用户信息并添加
del user with [xxx] # 删除名为xxx的用户
query info [*|all] # 查询所有的用户信息
query user by [yyy] # 根据yyy条件筛选用户
dis xxx # 禁用xxx用户
ena xxx # 启用xxx用户
unlock xxx # 解锁xxx用户
query user where day>1 # 查询未登入天数大于1的用户

###### ZABBIX ######
list [host|group] # 列出所有的主机|主机群组
query [host|group] by xxx # 列出以xxx关键字为条件的host|group
cfg export [a.json|a.xml] # 导出所有主机到文件中
create host form xxx.csv # 从xxx.csv穿件主机
```
# TODO LIST
[ ] 从文件中获取主机信息并创建