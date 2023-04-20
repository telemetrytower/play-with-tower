package main

import (
	"crawshaw.io/sqlite/sqlitex"
	"fmt"
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
	"github.com/julienschmidt/httprouter"
	"github.com/mergestat/go-mysql-sqlite-server/sqlitedb"
	"github.com/mergestat/go-mysql-sqlite-server/sqlitedb/tabledefine"
	"net/http"
	"os"
)

const (
	dbName = "testdata/Chinook_Sqlite.sqlite"
)

// create 5 promql http client
func main() {
	// 获取命令行参数
	args := os.Args
	//get dbname
	dbname:=args[1]
	fmt.Println(dbname)

	// add normal http serrver
	router := httprouter.New()
	PromService, err := tabledefine.Newtabledefine()
	if err == nil {
		router.GET("/v1/tableflag", PromService.Tableflags)
	} else {
		//log.Error("Failed to start reInference service:%v", err)
		panic(err)
	}

	go func() {
		if err := http.ListenAndServe(":3400", router); err != nil {
			//if err := http.ListenAndServe("10.19.34.183:8899", nil); err != nil {
			// TODO need check return value
			//log.Warn("listen up err:%v", err)
		}
		//logger.Log("listen up ok!")
	}()



	pool, err := sqlitex.Open(dbname, 0, 10)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := pool.Close(); err != nil {
			panic(err)
		}
	}()


	db := sqlitedb.NewDatabase(dbname, pool, nil,PromService)
	engine := sqle.NewDefault(
		sql.NewDatabaseProvider(
			db,
			information_schema.NewInformationSchemaDatabase(),
		))

	config := server.Config{
		Protocol: "tcp",
		Address:  "localhost:3306",
		Auth:     auth.NewNativeSingle("root", "", auth.AllPermissions),
	}

	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}




	s.Start()




}
