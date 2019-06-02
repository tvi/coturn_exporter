package parser

import (
	"testing"
)

func TestParse(t *testing.T) {
	body := `
<!DOCTYPE html>
<html>
  <head>
    <title>TURN Server (https admin connection)</title>
 <style> table, th, td { border: 1px solid black; border-collapse: collapse; text-align: left; padding: 5px;} table#msg th { color: red; background-color: white; } </style> </head>
  <body>
    <b>TURN Server</b><br><i>https admin connection</i><br>
 admin user: <b><i>user</i></b><br>
<br><a href="/home?realm=">home page</a><br>
<br><a href="/logout">logout</a><br>
<br>
<form action="/ps" method="POST">
  <fieldset><legend>Filter:</legend>
  <br>Realm name: <input type="text" name="realm" value="">  Client protocol: <input type="text" name="cprotocol" value="">  User name contains: <input type="text" name="puser" value=""><br><br>  Max number of output sessions in the page: <input type="text" name="maxsess" value="256"><br><br><input type="submit" value="Filter"></fieldset>
</form>
<br><b>TURN Sessions:</b><br><br><table>
<tr><th>N</th><th>Session ID</th><th>User</th><th>Realm</th><th>Origin</th><th>Age, secs</th><th>Expires, secs</th><th>Client protocol</th><th>Relay protocol</th><th>Client addr</th><th>Server addr</th><th>Relay addr (IPv4)</th><th>Relay addr (IPv6)</th><th>Fingerprints</th><th>Mobile</th><th>TLS method</th><th>TLS cipher</th><th>BPS (allocated)</th><th>Packets</th><th>Rate</th><th>Peers</th></tr>
<tr><td>1</td><td>000000000000015303<br><a href="/ps?cs=000000000000015303">cancel</a></td><td>user</td><td>north.gov</td><td></td><td>115</td><td>3485</td><td>TLS/TCP</td><td>UDP</td><td>8.8.8.8:36489</td><td>10.0.0.1:443</td><td>10.0.0.1:55046</td><td></td><td>ON</td><td>OFF</td><td>TLSv1.2</td><td>ECDHE-RSA-AES256-GCM-SHA384</td><td>0</td><td>rp=4060, rb=454992, sp=36, sb=4664
</td><td>r=5617, s=57, total=5674 (bytes per sec)
</td><td> 10.1.0.1 </td>
</table>
<br>Total sessions = 1<br>
</body>
</html>
	`
	i, _, err := Parse([]byte(body))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if i != 1 {
		t.Errorf("unexpected total sessions: got %v , want %v ", i, 1)
	}
}
