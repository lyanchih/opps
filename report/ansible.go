package report

import (
	"bytes"
	"github.com/spf13/cobra"
	"gitlab.adlinktech.com/lyan.hung/opps/conf"
	"log"
	"text/template"
)

const ansibleHostTmpl = `
{{- range $l, $nodes := .LabelMap }}
  {{- if eq $l "__DEFAULT__" }}
[DEFAULT]
  {{- else }}
[{{ $l }}]
  {{- end }}
  {{- range $_, $n := $nodes }}
{{$n.MAC}} ansible_host={{$n.IP}}
{{- if ne $.User ""}} ansible_user={{$.User}}{{- end }}
{{- if ne $.Passwd ""}} ansible_ssh_pass={{$.Passwd}}{{- end }} ansible_ssh_common_args='-o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no'
  {{- end }}
{{ end }}`

func init() {
	addReporter("ansible", &ansibleReporter{})
}

type sshData struct {
	User       string
	Passwd     string
	RootPasswd string
}

type ansibleData struct {
	sshData
	LabelMap map[string][]*conf.Node
}

type ansibleReporter struct {
	sshData
}

func (r *ansibleReporter) cmd(runE func(*cobra.Command, []string) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ansible",
		Short: "OPPS report ansible",
		RunE:  runE,
	}

	f := cmd.Flags()
	f.StringVarP(&r.User, "ansible-user", "", "", "Ansible ssh user")
	f.StringVarP(&r.Passwd, "ansible-passwd", "", "", "Ansible ssh passwd")
	f.StringVarP(&r.RootPasswd, "ansible-root-passwd", "", "", "Ansible ssh root user passwd")
	return cmd
}

func (r *ansibleReporter) Report(ss []*conf.Scenario) ([]byte, error) {
	data := ansibleData{
		sshData: r.sshData,
		LabelMap: map[string][]*conf.Node{
			"__DEFAULT__": []*conf.Node{},
		},
	}
	labelMap := data.LabelMap

	for _, s := range ss {
		for _, n := range s.Nodes {
			if len(n.IP) == 0 || len(n.MAC) == 0 {
				continue
			}

			for _, l := range n.Label {
				nodes, ok := labelMap[l]
				if ok {
					labelMap[l] = append(nodes, n)
				} else {
					labelMap[l] = []*conf.Node{n}
				}
				labelMap["__DEFAULT__"] = append(labelMap["__DEFAULT__"], n)
			}
		}
	}

	tmpl, err := template.New("ansible_hosts").Parse(ansibleHostTmpl)
	if err != nil {
		log.Println("Parse ansible_hosts template failed:", err)
		return nil, err
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		log.Println("Execute ansible_hosts template failed:", err)
		return nil, err
	}
	return buf.Bytes(), nil
}
