package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"security/models/ipModel"
	"security/models/udid"
	_ "security/routers"
	"security/service/lua"
	"security/service/timer"
	"security/service/timerUdid"
	"strings"
	"testing"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	. "github.com/smartystreets/goconvey/convey"
)

type argAny []interface{}

// get interface by index from interface slice
func (a argAny) Get(i int, args ...interface{}) (r interface{}) {
	if i >= 0 && i < len(a) {
		r = a[i]
	}
	if len(args) > 0 {
		r = args[0]
	}
	return
}

func ToStr(value interface{}) string {
	return fmt.Sprintf("%s", value)
}

func ValuesCompare(is bool, a interface{}, args ...interface{}) (ok bool, err error) {
	if len(args) == 0 {
		return false, fmt.Errorf("miss args")
	}
	b := args[0]
	arg := argAny(args)

	switch v := a.(type) {
	case reflect.Kind:
		ok = reflect.ValueOf(b).Kind() == v
	case time.Time:
		if v2, vo := b.(time.Time); vo {
			if arg.Get(1) != nil {
				format := ToStr(arg.Get(1))
				a = v.Format(format)
				b = v2.Format(format)
				ok = a == b
			} else {
				err = fmt.Errorf("compare datetime miss format")
				goto wrongArg
			}
		}
	default:
		ok = ToStr(a) == ToStr(b)
	}
	ok = is && ok || !is && !ok
	if !ok {
		if is {
			err = fmt.Errorf("expected: `%v`, get `%v`", b, a)
		} else {
			err = fmt.Errorf("expected: `%v`, get `%v`", b, a)
		}
	}

wrongArg:
	if err != nil {
		return false, err
	}

	return true, nil
}

func AssertIs(t *testing.T, a interface{}, args ...interface{}) {
	if ok, err := ValuesCompare(true, a, args...); ok == false {
		throwFail(t, err)
	}
}

func AssertNot(t *testing.T, a interface{}, args ...interface{}) {
	if ok, err := ValuesCompare(false, a, args...); ok == false {
		throwFail(t, err)
	}
}

func AssertString(t *testing.T, a interface{}, args ...interface{}) {

}

//判断是否是Int
func AssertInt(t *testing.T, data interface{}, args ...interface{}) {
	fmt.Println(data)
	intData := data.(int)
	fmt.Println(intData)
}

func getCaller(skip int) string {
	pc, file, line, _ := runtime.Caller(skip)
	fun := runtime.FuncForPC(pc)
	_, fn := filepath.Split(file)
	data, err := ioutil.ReadFile(file)
	var codes []string
	if err == nil {
		lines := bytes.Split(data, []byte{'\n'})
		n := 10
		for i := 0; i < n; i++ {
			o := line - n
			if o < 0 {
				continue
			}
			cur := o + i + 1
			flag := "  "
			if cur == line {
				flag = ">>"
			}
			code := fmt.Sprintf(" %s %5d:   %s", flag, cur, strings.Replace(string(lines[o+i]), "\t", "    ", -1))
			if code != "" {
				codes = append(codes, code)
			}
		}
	}
	funName := fun.Name()
	if i := strings.LastIndex(funName, "."); i > -1 {
		funName = funName[i+1:]
	}
	return fmt.Sprintf("%s:%d: \n%s", fn, line, strings.Join(codes, "\n"))
}

func throwFail(t *testing.T, err error, args ...interface{}) {
	if err != nil {
		con := fmt.Sprintf("\t\nError: %s\n%s\n", err.Error(), getCaller(3))
		if len(args) > 0 {
			parts := make([]string, 0, len(args))
			for _, arg := range args {
				parts = append(parts, fmt.Sprintf("%v", arg))
			}
			con += " " + strings.Join(parts, ", ")
		}
		t.Error(con)
		t.Fail()
	}
}

func throwFailNow(t *testing.T, err error, args ...interface{}) {
	if err != nil {
		con := fmt.Sprintf("\t\nError: %s\n%s\n", err.Error(), getCaller(2))
		if len(args) > 0 {
			parts := make([]string, 0, len(args))
			for _, arg := range args {
				parts = append(parts, fmt.Sprintf("%v", arg))
			}
			con += " " + strings.Join(parts, ", ")
		}
		t.Error(con)
		t.FailNow()
	}
}

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestMain is a sample to run an endpoint test
func TestMain(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestMain", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

/*
func TestRuleUpdate() {
	data := orm.Params{
		"reason": "检查到了大批的刷新2323",
	}
	rule.Update(&data, 1)
}
*/

func TestModelIpAdd(t *testing.T) {
	ipModel.Add(make([]string, 0), "", "", 0)
	var ipList = make([]string, 1)
	ipList[0] = "127.12.11.22"
	AssertIs(t, ipModel.Add(ipList, "desc", "unasm", 60), true)
}

// 需要测试三种情况，一种是什么都没有，一个是只有一个，一个是有多个的情况
func TestGetMin(t *testing.T) {
	data := ipModel.GetMinTime()
	AssertIs(t, data.Id, int64(-1))
	var (
		id = int64(39)
	)
	params := orm.Params{
		"expire_time": "1970-01-01 08:00:01",
	}
	AssertIs(t, ipModel.Update(&params, id), int64(1))

	data = ipModel.GetMinTime()

	AssertIs(t, data.Expire_time, "1970-01-01 08:00:01")
	params = orm.Params{
		"expire_time": "0000-00-00 00:00:00",
	}
	ipModel.Update(&params, id)
}

// 需要测试三种情况，一种是什么都没有，一个是只有一个，一个是有多个的情况
func TestTimerProcess(t *testing.T) {
	//common.SetMock()
	var (
		model ipModel.Black
	)
	ipArr := []string{"128.0.0.1", "127.0.0.12", "192.168.1.1", "192.168.0.1"}
	result := []bool{true, true, true, false}
	sendArr := ipArr[0:(len(ipArr) - 1)]
	lua.AddIp(sendArr, lua.Type_ip)
	for k, value := range ipArr {
		if k != 2 {
			continue
		}
		model = ipModel.GetByIp(value)
		AssertIs(t, timer.Process(model), result[k])
	}
	lua.DeleteIps(ipArr, lua.Type_ip)
}

/*
func TestClear(t *testing.T) {
	timer.NewTime <- "127.0.0.1"
	go timer.FakeClear()
	timer.NewTime <- "127.0.0.12"
}
*/

// 需要测试三种情况，一种是什么都没有，一个是只有一个，一个是有多个的情况
func TestUdid(t *testing.T) {
	go timerUdid.ClearPassData()
	//common.SetMock()
	udids := []string{"121", "12331"}
	for _, val := range udids {
		rs, _ := udid.Adds([]string{val}, "测试", "unasm", 1)
		udid.UpdateByIps(&orm.Params{"status": udid.Status_Online}, []string{val})

		fmt.Println(rs)
		timerUdid.NewIp <- val
	}
	time.Sleep(time.Second * 4)
	/*
		for _, val := range udids {
			udid.DeleteByUdids([]string{val})
		}
	*/
}
