package cmd

var cmdMap = map[string]func(string, string){
	"list":  lscmd,
	"query": lscmd,
	"cfg":   expAnyHosts,
}
