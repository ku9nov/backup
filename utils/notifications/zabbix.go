package notify

import (
	zabbix "github.com/blacked/go-zabbix"
	"github.com/ku9nov/backup/configs"
	"github.com/sirupsen/logrus"
)

func ZabbixSender(cfgValues configs.Config) {
	var metrics []*zabbix.Metric
	metrics = append(metrics, zabbix.NewMetric(cfgValues.Default.Host, cfgValues.Zabbix.ZabbixKey, cfgValues.Zabbix.ZabbixValue))
	packet := zabbix.NewPacket(metrics)
	z := zabbix.NewSender(cfgValues.Zabbix.ZabbixUrl, cfgValues.Zabbix.ZabbixPort)
	res, err := z.Send(packet)
	if err != nil {
		logrus.Errorln("ERROR sending to zabbix trapper: ", err)
		return
	}
	logrus.Infoln("Successfully sended to zabbix trapper. \n", string(res))
}
