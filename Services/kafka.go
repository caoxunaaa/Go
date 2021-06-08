package Services

import (
	"SuperxonWebSite/Models/ModuleQaStatisticDisplay"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tal-tech/go-queue/kq"
	"time"
)

var EmailOfWaringInfoKqPusher *kq.Pusher

func KafkaInit() {
	EmailOfWaringInfoKqPusher = kq.NewPusher([]string{"172.20.3.19:9092", "172.20.3.19:9093", "172.20.3.19:9094"}, "EmailOfWarningInfo")
}

//第一段：8点30传送到Kafka中，信息是昨天17点到今天8点30的告警数据; 第二段：13点传送到Kafka中，信息是今天8点30到今天13点的告警数据; 第三段：17点传送到Kafka中，信息是今天13点到今天17点的告警数据
func KafkaPushModuleWarningInfoByClock(clock int) error {
	now := time.Now()
	var startTime string
	var endTime string
	if clock == 8 {
		startTime = time.Date(now.Year(), now.Month(), now.Day()-1, 17, 0, 0, 0, now.Location()).Format("2006-01-02 15:04:05")
		endTime = time.Date(now.Year(), now.Month(), now.Day(), 8, 30, 0, 0, now.Location()).Format("2006-01-02 15:04:05")
	} else if clock == 13 {
		startTime = time.Date(now.Year(), now.Month(), now.Day(), 8, 30, 0, 0, now.Location()).Format("2006-01-02 15:04:05")
		endTime = time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, now.Location()).Format("2006-01-02 15:04:05")
	} else if clock == 17 {
		startTime = time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, now.Location()).Format("2006-01-02 15:04:05")
		endTime = time.Date(now.Year(), now.Month(), now.Day(), 17, 0, 0, 0, now.Location()).Format("2006-01-02 15:04:05")
	}

	var err error
	m := make(map[string][]struct {
		ModuleQaStatisticDisplay.QaStatisticInfo
		YellowLine int
		RedLine    int
	})
	m["ModuleNormalWarningInfo"], err = ModuleQaStatisticDisplay.GetAllModuleWaringInfo("TRX正常品", "%%", startTime, endTime)
	if err != nil {
		return err
	}
	m["ModuleRepairWarningInfo"], err = ModuleQaStatisticDisplay.GetAllModuleWaringInfo("TRX改制返工品", "%%", startTime, endTime)
	if err != nil {
		return err
	}
	fmt.Println("m", m)
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = EmailOfWaringInfoKqPusher.Push(string(data))
	if err != nil {
		return errors.New("kafka push error")
	}
	return nil
}
