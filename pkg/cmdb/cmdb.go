package cmdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"wagent/pkg/rpc"
	"wagent/pkg/wnet"

	"github.com/astaxie/beego/validation"
	"github.com/bitly/go-simplejson"
	"github.com/pbnjay/memory"
)

type Host struct {
	Hostname    string `json:"hostname" valid:"MaxSize(100)"`
	OS          string `json:"os" valid:"MaxSize(20)"`
	GID         uint   `json:"gid"`
	IP          string `json:"ip"`
	MemorySize  uint64 `json:"memory_size"`
	Cores       int    `json:"cores" valid:"Range(0,1000)"`
	Status      int    `json:"status" valid:"Range(0,1)"`
	Extras      string `json:"extras" valid:"MaxSize(3000)"`
	Uptime      int    `json:"uptime"`
	UpdatedBy   int    `json:"updated_by"`
	CreatedBy   int    `json:"created_by"`
	Description string `json:"description"`
}

type Guarder struct {
	Datacenter string `json:"datacenter"`
	Provider   string `json:"provider"`
	IP         string `json:"ip"`
	Port       uint   `json:"port"`
	Extras     string `json:"extras"`
}

type Likely struct {
	Hostname string `json:"hostname" valid:"MaxSize(100)"`
	IP       string `json:"ip"`
	Extras   string `json:"extras" valid:"MaxSize(3000)"`
	Uptime   int    `json:"uptime"`
	Status   int    `json:"status" valid:"Range(0,1)"`
}
type Unlikely struct {
	OS         string `json:"os" valid:"MaxSize(20)"`
	GID        uint   `json:"gid"`
	MemorySize uint64 `json:"memory_size"`
	Cores      int    `json:"cores" valid:"Range(0,1000)"`
}
type Ret struct {
	data string `json:"data"`
	msg  string `json:"msg"`
	code int    `json:"code"`
}

type GuarderRet struct {
	data []*Guarder
	msg  string `json:"msg"`
	code int    `json:"code"`
}
type HttpArgs struct {
	Method string
	Data   string
	Url    string
}

type HttpReply struct {
	Result string
	Error  string
	Code   int
}

const (
	CMDB_URL = "http://op.liuliancao.com"
)

func (h *Host) GetMemorySize() uint64 {
	return memory.TotalMemory()
}
func (h *Host) GetIP() []string {
	ips := wnet.GetAllIPs()
	return ips
}

func (h *Host) GetCores() int {
	return runtime.NumCPU()
}
func (h *Host) RegisterHost() int {
	var extrastr string
	extras, err := json.Marshal(h.GetExtras())
	if err != nil {
		extrastr = ""
	} else {
		extrastr = string(extras[:])
	}

	ips := h.GetIP()
	ipstr := ""
	for _, ip := range ips {
		ipstr = ipstr + ip + ","
	}

	host := &Host{
		Hostname:    h.GetHostname(),
		OS:          h.GetOS(),
		GID:         h.GetGID(),
		Cores:       h.GetCores(),
		MemorySize:  h.GetMemorySize(),
		IP:          ipstr,
		Extras:      extrastr,
		Uptime:      int(h.GetUptime()),
		Status:      1,
		CreatedBy:   2,
		UpdatedBy:   2,
		Description: "",
	}

	if err != nil {
		log.Println("parse host json failed", h, err)
	}
	ht, err := json.Marshal(host)
	fmt.Println(string(ht))
	result, errs := rpc.CMDBRpc("Post", string(ht), CMDB_URL+"/api/v1/host")
	if errs != "" {
		log.Fatalf("add host failed: %s", err.Error())
		return 0
	}
	var ret map[string]interface{}
	err = json.Unmarshal([]byte(result), &ret)
	if err != nil {
		log.Fatalf("parse add host json failed %s", err.Error())
		return 0
	}
	hid := int(ret["data"].(float64))
	if err != nil {
		log.Fatalln("get ret data id failed, ", err.Error())
		return 0
	}
	err = ioutil.WriteFile("./sn_id", []byte(strconv.Itoa(hid)), 0644)
	if err != nil {
		log.Fatalln("write sn_id failed", hid, err)
		return 0
	}
	log.Println("add host successfully!", host, hid)
	return hid

}
func (h *Host) MustExistsID() int {

	f, err := os.Open("./sn_id")
	if err != nil {
		log.Println(h.Hostname, "first init, registing to cmdb...", err)
		return h.RegisterHost()

	} else {
		hid, err := ioutil.ReadAll(f)

		if err != nil {
			log.Println("read sn_id caught error: ", err.Error())
			return 0
		}
		id, err := strconv.Atoi(string(hid))
		if id > 0 && err == nil {
			// check if sn_id exists in cmdb, if not need generate new
			result, _ := rpc.CMDBRpc("Get", "", CMDB_URL+"/api/v1/host/"+string(hid))

			var ret map[string]interface{}
			err = json.Unmarshal([]byte(result), &ret)
			if err != nil {
				log.Fatalf("parse query host json failed %s, %s", err.Error(), result)
				return 0
			}
			code := int(ret["code"].(float64))
			if code == 70004 {
				h.RegisterHost()
			} else {
				return id
			}

		}
	}
	return 0
}
func (h *Host) UploadHostInfoLikely() {
	var extrastr string
	extras, err := json.Marshal(h.GetExtras())
	if err != nil {
		extrastr = ""
	} else {
		extrastr = string(extras[:])
	}
	ips := h.GetIP()
	ipstr := ""
	for _, ip := range ips {
		ipstr = ipstr + ip + " "
	}

	//uptime := strconv.FormatInt(h.GetUptime(), 10)
	uptime := h.GetUptime()

	host := &Likely{
		Hostname: h.GetHostname(),
		IP:       ipstr,
		Extras:   extrastr,
		Uptime:   int(uptime),
		Status:   1,
	}
	log.Println(host)
	valid := validation.Validation{}
	ok, _ := valid.Valid(host)
	if !ok {
		log.Println(valid.Errors)
		return
	}

	f, err := os.Open("./sn_id")
	if err != nil {
		if err != nil {
			log.Panicln("get sn_id failed: ", err.Error())
			return
		}
	} else {
		hid, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatalf("read sn_id %s failed with %s", hid, err.Error())
			return
		}
		log.Println("update data is", host)
		ht, err := json.Marshal(host)
		fmt.Println(string(ht))
		result, errs := rpc.CMDBRpc("Put", string(ht), CMDB_URL+"/api/v1/host/"+string(hid))

		if errs != "" {
			log.Fatalf("update host %v failed %s", host, errs)
		}
		log.Println("update host successfully ", host, result)
	}

}
func (h *Host) UploadHostInfoUnLikely() {
	host := &Unlikely{
		OS:         h.GetOS(),
		GID:        h.GetGID(),
		Cores:      h.GetCores(),
		MemorySize: h.GetMemorySize(),
	}

	valid := validation.Validation{}
	ok, _ := valid.Valid(host)
	if !ok {
		log.Println(valid.Errors)
		return
	}

	f, err := os.Open("./sn_id")
	if err != nil {
		log.Panicln("get sn_id failed: ", err.Error())
		return

	} else {
		hid, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatalf("read sn_id %s failed with %s", hid, err.Error())
			return
		}
		ht, err := json.Marshal(host)
		result, errs := rpc.CMDBRpc("Put", string(ht), CMDB_URL+"/api/v1/host/"+string(hid))

		if errs != "" {
			log.Fatalf("update host", host, "failed", err.Error())
		}
		log.Println("update host successfully", host, result)
	}

}

func (g *Guarder) GetAllGuarders() string {
	result, errs := rpc.CMDBRpc("Get", "", CMDB_URL+"/api/v1/guarders")

	if errs != "" {
		log.Fatalf("get guarders failed %s", errs)
	}
	ret, _ := simplejson.NewJson([]byte(result))
	rows, err := ret.Get("data").Get("guarders").Array()
	if err != nil {
		log.Println("parse all guarders with err: ", err)
	}
	guardersStr := ""
	for _, row := range rows {
		if guarder, ok := row.(map[string]interface{}); ok {
			guardersStr += guarder["ip"].(string) + ":" + fmt.Sprintf("%s", guarder["port"]) + ","
		}
	}

	return strings.Trim(guardersStr, ",")
}

func (g *Guarder) SyncGuarders() {
	guarders := g.GetAllGuarders()
	err := ioutil.WriteFile("./guarders", []byte(guarders), 0644)
	if err != nil {
		log.Fatalln("write guarders failed", err.Error())
	}
}

func (g *Guarder) GetAvailableGuarders() {
	f, err := os.Open("./guarders")
	if err != nil {
		log.Println("open guarders caught error:", err.Error())
		g.SyncGuarders()
	}

	guarders, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println("read guarders caught error:", err.Error())
	}

	f, err = os.Open("./guarder")
	if err == nil {
		guarder, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println("cannot read guarder, maybe permission error ", err.Error())
		}
		if status, _ := rpc.HealthRpc(string(guarder)); status == 1 {
			return
		}
	}
	log.Println("open guarder error (ignore while first initting...)", err.Error())
	for _, guarder := range strings.Split(string(guarders), ",") {
		if status, _ := rpc.HealthRpc(guarder); status == 1 {
			ioutil.WriteFile("./guarder", []byte(guarder), 0644)
		}
	}
}
