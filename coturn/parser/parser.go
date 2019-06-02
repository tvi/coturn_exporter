package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tvi/coturn_exporter/coturn/types"
)

func parseLine(in []string) (types.Session, error) {
	ret := types.Session{}

	id, err := strconv.Atoi(strings.TrimSpace(in[0]))
	if err != nil {
		return ret, err
	}
	ret.Idx = id

	ret.ID = strings.TrimSpace(in[1])
	ret.User = strings.TrimSpace(in[3])
	ret.Realm = strings.TrimSpace(in[4])
	ret.Origin = strings.TrimSpace(in[5])

	age, err := strconv.Atoi(strings.TrimSpace(in[6]))
	if err != nil {
		return ret, err
	}
	ret.Age = age

	exp, err := strconv.Atoi(strings.TrimSpace(in[7]))
	if err != nil {
		return ret, err
	}
	ret.Expires = exp

	ret.ClientProto = strings.TrimSpace(in[8])
	ret.RelayProto = strings.TrimSpace(in[9])
	ret.ClientAddr = strings.TrimSpace(in[10])
	ret.ServerAddr = strings.TrimSpace(in[11])
	ret.RelayAddrV4 = strings.TrimSpace(in[12])
	ret.RelayAddrV6 = strings.TrimSpace(in[13])

	if strings.TrimSpace(in[14]) == "ON" {
		ret.Fingerprints = true
	}

	if strings.TrimSpace(in[15]) == "ON" {
		ret.Mobile = true
	}
	ret.TLSVers = strings.TrimSpace(in[16])
	ret.TLSCipher = strings.TrimSpace(in[17])

	bps, err := strconv.Atoi(strings.TrimSpace(in[18]))
	if err != nil {
		return ret, err
	}
	ret.BPS = bps

	pack, err := parsePackets(in[19])
	if err != nil {
		return ret, err
	}
	ret.Packets = pack

	rate, err := parseRate(in[20])
	if err != nil {
		return ret, err
	}
	ret.Rate = rate

	ret.Peers = strings.TrimSpace(in[21])
	return ret, nil
}

func parsePackets(in string) (types.Packets, error) {
	ret := types.Packets{}
	_, err := fmt.Sscanf(in, "rp=%d, rb=%d, sp=%d, sb=%d", &ret.RecvPackets, &ret.RecvBytes, &ret.SentPackets, &ret.SentBytes)
	return ret, err
}

func parseRate(in string) (types.Rate, error) {
	ret := types.Rate{}
	_, err := fmt.Sscanf(in, "r=%d, s=%d, total=%d", &ret.RecvBytesPerSec, &ret.SentBytesPerSec, &ret.TotalBytesPerSec)
	return ret, err
}

func Parse(body []byte) (int, map[string]types.Session, error) {
	s := string(body)
	re := regexp.MustCompile(`<tr><td>(.*?)</td><td>(.*?)<br><a href="/ps\?cs=(.*?)">cancel</a></td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>([a-z0-9\=,\W]*?)</td><td>([a-z0-9\=,\W]*?)</td><td>(.*?)</td>`)
	out := re.FindAllStringSubmatch(s, -1)

	sessions := map[string]types.Session{}
	for _, line := range out {
		s, err := parseLine(line[1:])
		if err != nil {
			return 0, nil, err
		}
		sessions[s.ID] = s
	}

	re = regexp.MustCompile(`<br>Total sessions = (.*)<br>`)
	outs := re.FindStringSubmatch(s)
	if len(outs) < 2 {
		return 0, nil, fmt.Errorf("did not find total sessions")
	}
	total, err := strconv.Atoi(outs[1])
	if err != nil {
		return 0, nil, fmt.Errorf("did not find total sessions: %v", err)
	}
	return total, sessions, nil
}
