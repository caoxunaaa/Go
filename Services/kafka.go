package Services

import (
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"SuperxonWebSite/Models/ModuleStatisticDisplay"
	"encoding/json"
	"errors"
	"github.com/tal-tech/go-queue/kq"
	"time"
)

var EmailOfWaringInfoKqPusher *kq.Pusher
var EmailOfWaringInfoWithStationKqPusher *kq.Pusher
var EmailOfPtrChangeInfoKqPusher *kq.Pusher

func KafkaInit() {
	EmailOfWaringInfoKqPusher = kq.NewPusher([]string{"172.20.3.19:9092", "172.20.3.19:9093", "172.20.3.19:9094"}, "EmailOfWarningInfo")
	EmailOfWaringInfoWithStationKqPusher = kq.NewPusher([]string{"172.20.3.19:9092", "172.20.3.19:9093", "172.20.3.19:9094"}, "EmailOfWarningInfoWithStation")
	EmailOfPtrChangeInfoKqPusher = kq.NewPusher([]string{"172.20.3.19:9092", "172.20.3.19:9093", "172.20.3.19:9094"}, "EmailOfPtrChangedInfo")
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
		ModuleStatisticDisplay.QaStatisticInfo
		YellowLine int
		RedLine    int
	})

	m["ModuleNormalWarningInfo"], err = ModuleRunDisplay.GetModuleWaringInfo("正常品", "%%", startTime, endTime)
	if err != nil {
		return err
	}
	m["ModuleRepairWarningInfo"], err = ModuleRunDisplay.GetModuleWaringInfo("改制返工品", "%%", startTime, endTime)
	if err != nil {
		return err
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if len(data) <= 0 {
		return errors.New("没有良率告警信息")
	}
	err = EmailOfWaringInfoKqPusher.Push(string(data))
	if err != nil {
		return errors.New("kafka push error")
	}
	return nil
}

func KafkaPushModuleWarningInfoWithStationByClock(clock int) error {
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
		ModuleRunDisplay.StationStatus
		YellowLine int
		RedLine    int
	})

	m["ModuleNormalWarningInfoWithStation"], err = ModuleRunDisplay.GetModuleWaringInfoWithStation("正常品", "%%", startTime, endTime)
	if err != nil {
		return err
	}
	m["ModuleRepairWarningInfoWithStation"], err = ModuleRunDisplay.GetModuleWaringInfoWithStation("改制返工品", "%%", startTime, endTime)
	if err != nil {
		return err
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if len(data) <= 0 {
		return errors.New("没有工位良率告警信息")
	}
	err = EmailOfWaringInfoWithStationKqPusher.Push(string(data))
	if err != nil {
		return errors.New("kafka push error")
	}
	return nil
}
