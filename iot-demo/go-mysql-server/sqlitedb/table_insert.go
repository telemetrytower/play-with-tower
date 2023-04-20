package sqlitedb

import (
	"fmt"
	"github.com/mergestat/go-mysql-sqlite-server/sqlitedb/promewrite"
	"strconv"
	"strings"

	"github.com/dolthub/go-mysql-server/sql"
)

var _ sql.InsertableTable = (*Table)(nil)
var _ sql.RowInserter = (*rowInserter)(nil)

// every table must assign a key
type keyflag struct{
	tablename string
	keys string
	timestampkey string
	keytime int64
}
var keyflagtest  keyflag


type rowInserter struct {
	*tableEditor
	table *Table
	proms *promewrite.HttpClient
}

func newRowInserter(table *Table) *rowInserter {
	return &rowInserter{
		tableEditor: newTableEditor(),
		table:       table,
	}
}

func (t *Table) Inserter(*sql.Context) sql.RowInserter {
	return newRowInserter(t)
}

func (inserter *rowInserter) Insert(ctx *sql.Context, row sql.Row) error {
	if inserter.table.db.options.PreventInserts {
		return ErrNoInsertsAllowed
	}
	conn := inserter.table.db.pool.Get(ctx)
	if conn == nil {
		return ErrNoSQLiteConn
	}
	defer inserter.table.db.pool.Put(conn)

	schema := inserter.table.Schema()
	colNames := make([]string, len(schema))
	colNames1 := make([]string, len(schema))
	colType := make([]interface{}, len(schema))
	for c, col := range schema {
		colNames[c] = fmt.Sprintf("'%s'", col.Name)
		colType[c] = fmt.Sprintf("%s", col.Type)
		colNames1[c] = fmt.Sprintf("%s", col.Name)
		//fmt.Print("##test,col.Name:",col.Name)
		//fmt.Print("##test,col.TYPE:",col.Type)
		/*if true == col.PrimaryKey {
			fmt.Print("##test,col PrimaryKey:",col.Name)
			keyflagtest.keys += col.Name
		}
		fmt.Print("##test,col.Comment:",col.Comment)
		if "SeriesTime" == col.Comment {
			keyflagtest.timestampkey = col.Name
		}*/
	}
	//fmt.Print("##test,tablename:",inserter.table.name)
	//fmt.Print("##test,colNames:",colNames)
	//fmt.Print("##test,row:",row)
	/*sql, args, err := sq.Insert(inserter.table.name).
		Columns(colNames...).
		Values(row...).
		ToSql()
	if err != nil {
		return err
	}

	err = sqlitex.Exec(conn, sql, nil, args...)
	if err != nil {
		return err
	}*/

	/*keyflagtest.tablename=inserter.table.name
	keyflagtest.keys += "usrid,"
	keyflagtest.keys += "sportmode,"*/
	//keyflagtest.timestampkey = "starttime"
	var metrics =make([]promewrite.MetricPoint,len(colNames))
	//fmt.Println("len(v.ColumnName):",len(colNames))
	var TagsMaptmp =make(map[string]string,len(colNames))
	var Metricname []string
	var colvalue [] float64

	proms:=inserter.table.db.promqlserver()
	fmt.Println("### test, insert proms:",proms)
	//fmt.Println("### test, inserter.table.name:",inserter.table.name)
	Tableflaginfo:=proms.Tableinfo[inserter.table.name]
	for i,temcol:=range colNames1{
		if Tableflaginfo==nil {
			continue
		}
		//fmt.Println("###test,temcol:",temcol)
		//fmt.Println("###test,keyflagtest.timestampkey:",keyflagtest.timestampkey)
		//if key=starttime, continue
		if strings.Contains(Tableflaginfo.Keytime,temcol) {
			//fmt.Println("###test,temcol0001:",temcol)
			//fmt.Print("##test,row type:",row[i])
			//fmt.Println("##test,row data:",row[i])
			keyflagtest.keytime = int64(row[i].(int32))
			continue
		}
		//get keylabel, usrid, sportmode
		if strings.Contains(Tableflaginfo.Keyflags,temcol) {
			//fmt.Print("##test,row type:",row[i].type)
			//fmt.Println("##test,Keyflags:",temcol,row[i])
			/*if colType[i].(string) == "INT" {
				temp1,err:= strconv.ParseInt(row[i].(string), 10,64)
				if err != nil {
					// 可能字符串 s 不是合法的整数格式，处理错误
				} else {
					fmt.Print("##test temp1:",temp1)
					TagsMaptmp[temcol] = row[i].(string)
				}
			}*/

			if op,ok:=row[i].(string) ;ok {
				TagsMaptmp[temcol] = op
				//fmt.Println("##test,string,Keyflags:",temcol,row[i],op)
			}
			if op,ok:=row[i].(int32);ok {
				//fmt.Println("##test,int,Keyflags:",temcol,row[i],op)
				TagsMaptmp[temcol] = strconv.FormatInt(int64(op), 10)

			}
			//colvalue = append(colvalue,float64(v.ValueExpr[i].(int64)))
			continue
		}
		Metricname = append(Metricname,Tableflaginfo.Tablename+temcol)
		switch colType[i].(string) {
		case "INT":
			colvalue = append(colvalue,float64(row[i].(int32)))
		case "FLOAT":
			colvalue = append(colvalue,float64(row[i].(float32)))
		}

		//colvalue = append(colvalue,float64(row[i].(int)))
	}

	for i,tempmetric:=range Metricname {
		metrics[i] = promewrite.MetricPoint{
			Metric:  tempmetric,
			TagsMap: TagsMaptmp,
			Time:    keyflagtest.keytime,
			Value:   colvalue[i],
		}

		proms.MetricCh<-&metrics[i]
		//fmt.Println("##test,write to batch channel MetricCh:",i,len(tempmetric),tempmetric,metrics[i].TagsMap,metrics[i].Time)
	}

	//proms.MetricCh<-metrics

	// get promeql client
	/*inserter.table.db.proms.PromqServer.RemoteWrite(metrics)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("prometheus write end")*/

	return nil
}

func (inserter *rowInserter) Close(*sql.Context) error { return nil }
