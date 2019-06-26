package metric

import (
	"strconv"
	"testing"

	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/vitelabs/go-vite/log15"
)

func TestCounter(t2 *testing.T) {

	r := metrics.NewRegistry()

	c := metrics.NewCounter()
	r.Register("foo", c)

	c.Inc(1)

	r.Each(func(k string, c interface{}) {
		println("key:" + k + ",v:" + strconv.FormatInt(c.(metrics.Counter).Count(), 10))
	})

	//m := metrics.NewMeter()
	//
	//go func() {
	//	for {
	//		m.Mark(1)
	//		time.Sleep(20 * time.Millisecond)
	//	}
	//}()
	//
	//go func() {
	//	t := time.NewTicker(time.Second * 2)
	//	for {
	//		select {
	//		case <-t.C:
	//			println(strconv.Itoa(int(m.Rate1())))
	//		}
	//	}
	//
	//}()

	t := metrics.NewTimer()

	testTime(t)
	println(t.Sum())
	println(t.Count())

}

func testTime(t metrics.Timer) {
	defer t.UpdateSince(time.Now())

	time.Sleep(2 * time.Second)
}

//设置日志文件地址
func SetTestLogContext() {
	dir := "/Users/jie/log/test.log"
	log15.Info(dir)
	log15.Root().SetHandler(
		log15.LvlFilterHandler(log15.LvlInfo, log15.Must.FileHandler(dir, log15.JsonFormat())),
	)

}

//打印日志
func TestPrint(t *testing.T) {
	SetTestLogContext()
	//INFO[09-07|15:09:44] msg:                                     type=1 appkey=govite group=msg name=effectivemsg metric=1 class=vm code=6001600201602080919052602090f3
	logger := log15.New("logtype", "1", "appkey", "govite", "group", "monitor", "class", "vm")

	logger.Info("name and type", "name", "pool", "metric", time.Now().Sub(time.Now().Add(time.Second*10)).Nanoseconds()/time.Millisecond.Nanoseconds())

}
