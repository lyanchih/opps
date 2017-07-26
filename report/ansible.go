package report

import (
	"bytes"
	"gitlab.adlinktech.com/lyan.hung/opps/conf"
	"log"
	"text/template"
)

const ansibleHostTmpl = `
{{- range $l, $nodes := . }}
  {{- if eq $l "__DEFAULT__" }}
[DEFAULT]
  {{- else }}
[{{ $l }}]
  {{- end }}
  {{- range $_, $n := $nodes }}
{{$n.MAC}} ansible_host={{$n.IP}} ansible_user=adlink ansible_ssh_pass=adlink ansible_ssh_common_args='-o StrictHostKeyChecking=no'
  {{- end }}
{{ end }}`

func init() {
	addReporter("ansible", ansibleReporter{})
}

type ansibleReporter struct{}

func (r ansibleReporter) Report(ss []*conf.Scenario) ([]byte, error) {
	labelMap := map[string][]*conf.Node{
		"__DEFAULT__": []*conf.Node{},
	}
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
	err = tmpl.Execute(buf, labelMap)
	if err != nil {
		log.Println("Execute ansible_hosts template failed:", err)
		return nil, err
	}
	return buf.Bytes(), nil
}
