package port

import (
	"sort"
	"strconv"
	"strings"
)

var sokets = map[string]string{
	"tcp":  "tcp_socket",
	"udp":  "udp_socket",
	"sctp": "sctp_socket",
}

func CreatePolicyFromInspect(ports []map[string]interface{}, strPolicy string) (string, error) {
	sort.Slice(ports, func(i, j int) bool {
		intPortNumber_i, err_i := strconv.Atoi(ports[i]["portNumber"].(string))
		intPortNumber_j, err_j := strconv.Atoi(ports[j]["portNumber"].(string))
		if err_i != nil || err_j != nil {
			return false
		}
		return intPortNumber_i < intPortNumber_j
	})
	for _, port := range ports {
		intPortNumber, err := strconv.Atoi(port["portNumber"].(string))
		if err != nil {
			return "", err
		}
		portContext, _ := listPorts(intPortNumber, port["portNumber"].(string))
		strPolicy += ("    (allow process %s ( %s ( name_bind )))\n" + portContext + sokets[port["protocol"].(string)])
	}
	return strPolicy, nil
}

// internal functions
func listPorts(portNumber int, portProto string) (string, error) {
	portContexts, err := getPortContexts()
	if err != nil {
		return "", err
	}

	for _, port := range portContexts {
		proto := port.proto
		low := port.low
		high := port.high
		if low <= portNumber && portNumber <= high && strings.EqualFold(proto, portProto) {
			return port.ctype, nil
		}
	}

	return "", nil
}

type portContext struct {
	ctype string
	proto string
	low   int
	high  int
}

func getPortContexts() ([]portContext, error) {
	return []portContext{
		{"http_port_t", "tcp", 80, 80},
		{"ssh_port_t", "tcp", 22, 22},
		{"dns_port_t", "udp", 53, 53},
	}, nil
}
