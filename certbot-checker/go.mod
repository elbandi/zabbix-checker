module certbot-checker

go 1.23.12

require (
	github.com/urfave/cli/v2 v2.27.7
	golang.zabbix.com/agent2 v0.0.0-20251031143913-c525822dce02
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20250705151800-55b8f293f342 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	golang.zabbix.com/sdk v1.2.2-0.20251003101739-7de6a6ab9506 // indirect
)

replace golang.zabbix.com/agent2 v0.0.0-20251031143913-c525822dce02 => github.com/zabbix/zabbix/src/go v0.0.0-20251031143913-c525822dce02
