package b1815

import "time"

func (client *Client) GetIdleTimes() (lastRecieve time.Time, logon time.Time) {
	return client.lastReceive, client.logonTime
}
