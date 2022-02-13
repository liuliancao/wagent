package wcron

import (
	"fmt"
	"wagent/pkg/cmdb"

	"github.com/robfig/cron"
)

func hello() {
	fmt.Println("hello")
}
func Run() {
	var host cmdb.Host
	//var guarder cmdb.Guarder
	c := cron.New()
	//second, minute, hour, day, month, weekday
	//c.AddFunc("*/20 * * * * *", host.UploadHostInfoLikely)
	c.AddFunc("*/20 * * * * *", host.UploadHostInfoUnLikely)
	//c.AddFunc("*/20 * * * * *", guarder.GetAvailableGuarders)
	//c.AddFunc("*/20 * * * * *", guarder.SyncGuarders)
	c.Start()
}
