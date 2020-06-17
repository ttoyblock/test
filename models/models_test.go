/*
 * @Author: ttoy
 * @Date: 2020-03-06 00:39:09
 * @Last Modified by:
 * @Last Modified time: 2020-03-06 00:47:09
 */
package models_test

import (
	"database/sql"
	"testing"
	"toolkit/models"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	path := "root:pass@tcp(alpha:3306)/ultrax?clientFoundRows=false&parseTime=true&loc=Asia%2FShanghai&timeout=5s&collation=utf8mb4_bin&interpolateParams=true"
	db, _ = sql.Open("mysql", path)
}

func Test_Get(t *testing.T) {
	m := new(postmodel.PreForumThread)
	m.Tid = 63328

	err := models.Get(db, m)
	if err != nil {
		t.Error(err)
	}
	t.Log(*m)
}

func Test_Gets(t *testing.T) {
	tids := []int64{63328, 503251}
	ms := make([]*postmodel.PreForumThread, 0)

	err := models.Gets(db, tids, postmodel.ModelHelperPreForumThread, &ms)
	if err != nil {
		t.Error(err)
	}
	for _, v := range ms {
		t.Log(*v)
	}
}

func Test_TxGet(t *testing.T) {
	tx, _ := db.Begin()
	defer tx.Commit()
	m := new(postmodel.PreForumThread)
	where := map[string]interface{}{"tid": 63328}

	err := models.TxGet(tx, m, where)
	if err != nil {
		t.Error(err)
	}
	t.Log(*m)
	tx.Commit()
}

func Test_DBGetOne(t *testing.T) {
	m := new(postmodel.PreForumThread)
	where := map[string]interface{}{"Tid": 63328}

	err := models.DBGetOne(db, m, where)
	if err != nil {
		t.Error(err)
	}
	t.Log(*m)
}

func Test_DBTxGetOne(t *testing.T) {
	m := new(postmodel.PreForumThread)
	where := map[string]interface{}{"Tid": 63328}
	tx, _ := db.Begin()
	defer tx.Commit()

	err := models.DBTxGetOne(tx, m, where)
	if err != nil {
		t.Error(err)
	}
	t.Log(*m)
}

func Test_DBGetMulti(t *testing.T) {
	tids := []int64{63328, 503251}
	where := map[string]interface{}{"tid in": tids}
	ms := make([]*postmodel.PreForumThread, 0)

	err := models.DBGetMulti(db, &ms, where)
	if err != nil {
		t.Error(err)
	}
	for _, v := range ms {
		t.Log(*v)
	}
}

func Test_DBInsert(t *testing.T) {
	m := new(postmodel.PreForumThread)
	data := []map[string]interface{}{
		map[string]interface{}{"fid": 28, "subject": "test subject1"},
		map[string]interface{}{"fid": 29, "subject": "test subject2"},
	}

	id, err := models.DBInsert(db, m, data)
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
}

func Test_DBTxInsert(t *testing.T) {
	tx, _ := db.Begin()
	defer tx.Commit()
	m := new(postmodel.PreForumThread)
	data := []map[string]interface{}{
		map[string]interface{}{"fid": 38, "subject": "test subject3"},
		map[string]interface{}{"fid": 39, "subject": "test subject4"},
	}

	id, err := models.DBTxInsert(tx, m, data)
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
}

func Test_DBUpdate(t *testing.T) {
	m := new(postmodel.PreForumThread)
	data := map[string]interface{}{"price": 28}
	where := map[string]interface{}{"tid": 63328}

	id, err := models.DBUpdate(db, m, where, data)
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
}

func Test_DBTxUpdate(t *testing.T) {
	tx, _ := db.Begin()
	defer tx.Commit()

	m := new(postmodel.PreForumThread)
	data := map[string]interface{}{"price": 38}
	where := map[string]interface{}{"tid": 63328}

	id, err := models.DBTxUpdate(tx, m, where, data)
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
}

func Test_DBDelete(t *testing.T) {
	m := new(postmodel.PreForumThread)
	where := map[string]interface{}{"tid": 632224}

	aff, err := models.DBDelete(db, m, where)
	if err != nil {
		t.Error(err)
	}
	t.Log(aff)
}

func Test_DBTxDelete(t *testing.T) {
	tx, _ := db.Begin()
	defer tx.Commit()
	m := new(postmodel.PreForumThread)
	where := map[string]interface{}{"tid": 632225}

	aff, err := models.DBTxDelete(tx, m, where)
	if err != nil {
		t.Error(err)
	}
	t.Log(aff)
}
