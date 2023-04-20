package tabledefine

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/mergestat/go-mysql-sqlite-server/sqlitedb/promewrite"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Tableflaginfo struct {
	Tablename string `json:"tablename"`
	Keyflags string `json:"keyflags"`
	Keytime string `json:"keytime"`
}
type TableflaginfoEntries struct {
	TableflaginfoEntries []*Tableflaginfo  `json:"tableflaginfoEntries"`
}


type Promqldata struct {
	Tableinfo map[string] *Tableflaginfo
	PromqServer *promewrite.HttpClient
	MetricCh chan *promewrite.MetricPoint
	CloseCh chan bool
	TableflaginfoEntries []*Tableflaginfo  `json:"tableflaginfoEntries"`
}

//var tableinfo map[string] *Tableflaginfo
func initlocaltabelfile(pq *Promqldata){
	// check whether to create new data file
	var file *os.File
	var err error

	file,err=os.Open("table.json")
	if err!=nil {
		fmt.Println("fail to open new datafile")
		return
	}

	defer file.Close()

	//allprops := new(TableflaginfoEntries)
	/*br := bufio.NewReader(file)
	for {
		json_message,_, c := br.ReadLine() //按行读取文件
		if c == io.EOF {
			break
		}
		fmt.Println(" json_message:",json_message[:])
		rtep:=new (Tableflaginfo)
		err := json.Unmarshal(json_message[:len(json_message)], rtep)
		fmt.Println(" rtep:",rtep)
		if err != nil {
			fmt.Println(" cannot Unmarshal data file content")
		}

		pq.TableflaginfoEntries=append(pq.TableflaginfoEntries,rtep)
	}*/
	decoder  :=json.NewDecoder(file)
	for decoder.More() {
		var Tableflaginfotmp Tableflaginfo
		err:=decoder.Decode(&Tableflaginfotmp)
		if err != nil {
			fmt.Println(" cannot Decode data file content")
		}
		fmt.Println("Decode Tableflaginfotmp",Tableflaginfotmp)
		pq.TableflaginfoEntries=append(pq.TableflaginfoEntries,&Tableflaginfotmp)
	}
	fmt.Println("read table data file:",pq.TableflaginfoEntries)

}
func Newtabledefine() (*Promqldata, error)  {

	// backend
	backend:="https://io.telemetrytower.com/api/v1/push"  // cortex backend
	//backend:="http://10.19.34.132:9090/api/v1/write"  // prometheus backend
	promqServer, err := promewrite.NewClient(backend, 10*time.Second)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	pq := &Promqldata{
		Tableinfo: make(map[string] *Tableflaginfo),
		PromqServer: promqServer,
		MetricCh: make(chan *promewrite.MetricPoint, 1000),
		CloseCh: make(chan bool),
		TableflaginfoEntries: make([]*Tableflaginfo,10),
	}

	// init local table file
	initlocaltabelfile(pq)

	// init pq.Tableinfo
	lentmp:=len(pq.TableflaginfoEntries)
	fmt.Println("##test,lentmp:",lentmp)
	if lentmp > 0 {
		for _,table:=range pq.TableflaginfoEntries {
			if table == nil {
				continue
			}
			tablenametmp:=table.Tablename
			pq.Tableinfo[tablenametmp]=table
		}
	}


	go pq.BatchWrite()


	return pq,nil

}
func (pq *Promqldata) Tableflags(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	// get param
	m, _ := url.ParseQuery(r.URL.RawQuery)
	fmt.Println("origin:",m)
	t := Tableflaginfo{}
	// get tablename
	fmt.Println(m["table"])
	tablename, ok := m["table"]
	if !ok {
		tablename[0] = ""
	} else {
		t.Tablename = tablename[0]
	}

	// get keyflags
	keyflags, ok := m["keyflags"]
	if !ok {
		keyflags[0] = ""
	} else {
		for _,k:=range keyflags {
			t.Keyflags += k
			t.Keyflags += ","
		}
	}

	// get keyflags
	keytime, ok := m["keytime"]
	if !ok {
		keytime[0] = ""
	} else {
		t.Keytime = keytime[0]
	}

	// store table info
	pq.Tableinfo[tablename[0]]=&t
	fmt.Println("t:",t)
	fmt.Println(tablename[0])
	fmt.Println(keyflags)
	fmt.Println(keytime[0])
	//fmt.Println(tableinfo)
	msg, err := json.Marshal(t)
	if err !=nil {
		fmt.Println("json.Marshal:",err)
	}
	w.WriteHeader(200)
	w.Write(msg)

	// write data to local var
	pq.TableflaginfoEntries=append(pq.TableflaginfoEntries,&t)

	// write data to local file
	// check whether to create new data file
	var file *os.File
	//file,err=os.Open("table.json")
	file, err = os.OpenFile("table.json", os.O_RDWR|os.O_APPEND, os.ModeAppend)
	if err!=nil {
		file, err = os.Create("table.json")
		if err != nil {
			fmt.Println(" dont find data file,create new one", err)
		}
	}
	defer file.Close()
	jsonData, err := json.Marshal(t)
	if err!= nil {
		fmt.Println("fail to Marshal table data")
	}
	fmt.Println("waiting write json data:",jsonData)
	_,err =file.WriteString(string(jsonData) + "\n")
	if err!= nil {
		fmt.Println("fail to write table data to file!",err)
	}

}


func (pq *Promqldata) BatchWrite() {
	for {
		select {
		case <-pq.CloseCh:
			return
		case <-time.After(1 * time.Second):
			l := len(pq.MetricCh)
			if l == 0 {
				fmt.Println("MetricCh is empty,continue")
				continue
			}
			i := int(0)
			var metrics =make([]promewrite.MetricPoint,5)
			for ; i < l && i < 5; i++ {
				metricstmp:= <-pq.MetricCh
				metrics=append(metrics,*metricstmp)
			}
			err:=pq.PromqServer.RemoteWrite(metrics)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("prometheus write end")
		}
	}

}