package test

// package test

// import (
// 	"fmt"
// 	"strconv"
// 	"strings"
// 	"testing"

// 	"github.com/0x1un/omtools/parser/zbxcsv"

// 	"github.com/0x1un/go-zabbix"

// 	"github.com/antlr/antlr4/runtime/Go/antlr"
// )

// var InterfaceType = map[string]int{
// 	"agent": 1,
// 	"snmp":  2,
// 	"ipmi":  3,
// 	"jmx":   4,
// }

// type BaseZbxListener struct {
// 	zbxcsv.BasezbxcsvListener
// 	hosts    zabbix.CreateHostRequest
// 	interfac zabbix.Interface
// }

// func NewBaseZbxListener() *BaseZbxListener {
// 	return &BaseZbxListener{
// 		zbxcsv.BasezbxcsvListener{},
// 		zabbix.CreateHostRequest{},
// 		zabbix.Interface{},
// 	}
// }

// // ExitHost is called when production host is exited.
// func (s *BaseZbxListener) ExitHost(ctx *zbxcsv.HostContext) {
// 	s.hosts.Host = ctx.GetText()
// }

// // ExitPort is called when production port is exited.
// func (s *BaseZbxListener) ExitPort(ctx *zbxcsv.PortContext) {
// 	port, err := strconv.Atoi(ctx.GetText())
// 	if err != nil {
// 		panic(err)
// 	}
// 	s.interfac.Port = port
// }

// // ExitIp is called when production ip is exited.
// func (s *BaseZbxListener) ExitIp(ctx *zbxcsv.IpContext) {
// 	s.interfac.IP = ctx.GetText()
// }

// // ExitIface_type is called when production iface_type is exited.
// func (s *BaseZbxListener) ExitIface_type(ctx *zbxcsv.Iface_typeContext) {
// 	s.interfac.Type = InterfaceType[strings.ToLower(ctx.GetText())]
// }

// // ExitDns is called when production dns is exited.
// func (s *BaseZbxListener) ExitDns(ctx *zbxcsv.DnsContext) {
// 	if ss := ctx.GetText(); ss != "" {
// 		s.interfac.DNS = ss
// 	}
// }

// // ExitHostname is called when production hostname is exited.
// func (s *BaseZbxListener) ExitHostname(ctx *zbxcsv.HostnameContext) {
// 	s.hosts.Host = ctx.GetText()
// }

// // ExitVis_name is called when production vis_name is exited.
// func (s *BaseZbxListener) ExitVis_name(ctx *zbxcsv.Vis_nameContext) {
// 	s.hosts.VisibleName = ctx.GetText()
// }

// // ExitGroup_id is called when production group_id is exited.
// func (s *BaseZbxListener) ExitGroup_id(ctx *zbxcsv.Group_idContext) {
// 	s.hosts.Groups = append(s.hosts.Groups, zabbix.Group{GroupID: ctx.GetText()})
// }

// // ExitTempl_id is called when production templ_id is exited.
// func (s *BaseZbxListener) ExitTempl_id(ctx *zbxcsv.Templ_idContext) {

// }

// // ExitDisable_host is called when production disable_host is exited.
// func (s *BaseZbxListener) ExitDisable_host(ctx *zbxcsv.Disable_hostContext) {
// 	s.hosts.Status = func(ss string) int {
// 		if ss == "1" {
// 			return 1
// 		}
// 		return 0
// 	}(ctx.GetText())
// }

// // ExitUse_ip is called when production use_ip is exited.
// func (s *BaseZbxListener) ExitUse_ip(ctx *zbxcsv.Use_ipContext) {
// 	if ctx.GetText() == "1" {
// 		s.interfac.DNS = ""
// 		s.interfac.Useip = 1
// 	}
// }

// // ExitDescription is called when production description is exited.
// func (s *BaseZbxListener) ExitDescription(ctx *zbxcsv.DescriptionContext) {
// 	s.hosts.Description = ctx.GetText()
// }

// func (s *BaseZbxListener) VisitTerminal(node antlr.TerminalNode) {
// }

// func TestBasecommandsVisitor(t *testing.T) {
// 	input := "CD-GPA3FVK-SWiii-04,A3-接入交换机-4iii,54,172.19.2.10:10050+agent&172.19.2.11:10050+snmp&172.19.2.12:10050+ipmi,4321,1,1,8.8.8.8,测试的一号主机"
// 	fs := antlr.NewInputStream(input)
// 	lex := zbxcsv.NewzbxcsvLexer(fs)
// 	tokens := antlr.NewCommonTokenStream(lex, antlr.TokenDefaultChannel)
// 	p := zbxcsv.NewzbxcsvParser(tokens)
// 	p.BuildParseTrees = true
// 	tree := p.Row()
// 	listener := NewBaseZbxListener()
// 	antlr.ParseTreeWalkerDefault.Walk(listener, tree)
// 	fmt.Printf("%#v\n", listener.hosts)
// }
