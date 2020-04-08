/*
 * @Author: ttoy
 * @Date: 2020-03-04 23:07:07
 * @Last Modified by:
 * @Last Modified time: 2020-03-06 00:39:17
 */
package models

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"
	"toolkit/cache"

	"gitee.com/ikongjix/go_common/redis_db"
	"gitee.com/ikongjix/go_common/redis_db/master_db"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/go-redis/redis"
)

const (
	// PropsKeyFormat _
	PropsKeyFormat = "props:%v:%v" // tablename:id
)

// Model is a interface that abstract all behaviours of models
type Model interface {
	TableName() string
	PkName() string
	PkValue() int64
	CacheKey() string
	CacheExpireTime() time.Duration
	Flush() error
}

// type CacheWarmer interface {
// 	Model
// 	FormatTime()
// }

// Get is a function used to get data from model with cache
// REQ: instancePtr's id is required.
func Get(db *sql.DB, instancePtr Model) error {
	if propsModel, ok := instancePtr.(PropsModel); ok {
		err := LoadProps(propsModel)
		if err != nil && err != redis.Nil {
			return err
		}
	}

	cacheKey := instancePtr.CacheKey()
	if instancePtr.PkValue() > 0 && instancePtr.CacheKey() != "" {
		err := cache.GetCachedStruct(cacheKey, instancePtr)
		if err == nil {
			return nil
		} else if err != redis.Nil {
			return err
		}
	}

	err := DBGetOne(db, instancePtr, map[string]interface{}{instancePtr.PkName(): instancePtr.PkValue()})
	if err == sql.ErrNoRows {
		// NOTE:
		// 理想情况是，在查询不到 row 时，
		// 1. 返回错误
		// 2. 将 instance 置为 nil
		// 但是我太挫，没有实现，当前不返回 nil 而返回 对应 model 的 zero

		// 赋值成功，但是赋值结果为 model 的 zero 值，而非 nil
		v := reflect.ValueOf(instancePtr)
		v.Elem().Set(reflect.Zero(v.Elem().Type()))

		// 直接赋值失败 - 只是赋值给了 interface，没有打穿到 *model
		// instance = nil

		// 尝试将 *model 赋值为 nil, 未成功 -- 原因是，interface 的 elem 地址 是 unaddressable 的
		// 得到 instance 的 value -- *user.User
		// v := reflect.ValueOf(instance)

		// fmt.Println(v.CanSet())
		// vType := v.Type()
		// fmt.Println(vType)
		// ptr := reflect.PtrTo(v.Elem().Type())
		// fmt.Println(ptr)
		// val := reflect.Zero(ptr)
		// fmt.Println(val)
		// v.Set(val)

		return err
	}

	if err != nil {
		return err
	}

	if instancePtr.CacheKey() != "" {
		cache.CacheStruct(cacheKey, instancePtr, instancePtr.CacheExpireTime())
	}
	return nil
}

// TxGet get with props
func TxGet(tx *sql.Tx, instancePtr Model, where map[string]interface{}) error {
	err := DBTxGetOne(tx, instancePtr, where)
	if err == sql.ErrNoRows {
		v := reflect.ValueOf(instancePtr)
		v.Elem().Set(reflect.Zero(v.Elem().Type()))
		return sql.ErrNoRows
	}

	if err != nil {
		return err
	}

	if propsModel, ok := instancePtr.(PropsModel); ok {
		err = LoadProps(propsModel)
		if err != nil && err != redis.Nil {
			return err
		}
	}

	return nil
}

// Currently use this modelHelper function to convert all models to Model interface.
// May do some promotion on this in the future.
type modelHelper func([]int64) []Model

// Gets is a function used to get multi-data from model one by one
// NOTE:
// This function has three parameters:
// 1. ids -- the target ids of instances you desired to get.
// 2. helper -- a model package level function used to convert model instaces to Model interfaces
//              you have to implement this for for the model you defined.
// 3. instanceSlicePtr -- a pointer pointed to a slice for resultes
//
// Take user.User model as an example,
// we have helper named user.ModelInstanceHelper
// this method could be called as follow:
// target := make([]*user.User)
// models.Gets(db, []{1, 2, 3, 4}, user.ModelInstanceHelper, &target)
func Gets(db *sql.DB, ids []int64, helper modelHelper, instaceSlicePtr interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	instances := helper(ids)
	if instances[0].CacheKey() == "" {
		return errors.New("cache key is empty")
	}

	instancesAsInterfaces := make([]interface{}, len(instances))
	for i, v := range instances {
		instancesAsInterfaces[i] = v
	}

	keys := make([]string, len(instances))
	for i, v := range instances {
		keys[i] = v.CacheKey()
	}

	unCachedIndex, err := cache.MGetCachedStructs(keys, instancesAsInterfaces)
	if err != nil && err != redis.Nil {
		return err
	}
	if len(unCachedIndex) > 0 {
		for _, idx := range unCachedIndex {
			unCachedInstance := instances[idx]
			tmpWhere := map[string]interface{}{unCachedInstance.PkName(): unCachedInstance.PkValue()}
			err := DBGetOne(db, unCachedInstance, tmpWhere)
			if err == nil {
				instancesAsInterfaces[idx] = unCachedInstance
				unCachedInstanceKey := unCachedInstance.CacheKey()
				cache.CacheStruct(unCachedInstanceKey, unCachedInstance, redis_db.ONE_MONTH)
			}
		}
	}

	for _, v := range instancesAsInterfaces {
		if v != nil {
			if propsModel, ok := v.(PropsModel); ok {
				err = LoadProps(propsModel)
				if err != nil && err != redis.Nil {
					return err
				}
			}
			slicePointerValue := reflect.ValueOf(instaceSlicePtr).Elem()
			slicePointerValue.Set(reflect.Append(slicePointerValue, reflect.ValueOf(v)))
		}
	}

	return nil
}

// PropsModel is an interface abstracts the behaviours of model with propties
type PropsModel interface {
	Model
	EncodingProps() ([]byte, error)
	DecodingProps(propsValue []byte) error
}

// PropsKey is a function that used to get props key of model instance
func PropsKey(instance PropsModel) string {
	return fmt.Sprintf(PropsKeyFormat, instance.TableName(), instance.PkValue())
}

// LoadProps - function to load props for PropsModel instances
func LoadProps(instance PropsModel) error {
	propsValue, err := GetProps(instance)
	if err != nil {
		return err
	}

	err = instance.DecodingProps(propsValue)
	if err != nil {
		return err
	}

	return nil
}

// GetProps is a function used to get the props of a Model instance
func GetProps(instance PropsModel) (v []byte, err error) {
	key := PropsKey(instance)
	v, err = master_db.MasterRedis.Get(key).Bytes()
	return
}

// SaveProps is a function that used to save the props of a Model instance
func SaveProps(instance PropsModel) error {
	key := PropsKey(instance)
	value, err := instance.EncodingProps()
	if err != nil {
		return err
	}
	master_db.MasterRedis.Set(key, value, 0)
	// TODO: remove - make sure that the cache logic will not cache props info
	instance.Flush()
	return nil
}

// DeleteProps is a function used to delete props of a Model instance
func DeleteProps(instance PropsModel) error {
	key := PropsKey(instance)
	err := master_db.MasterRedis.Del(key).Err()
	// TODO: remove - make sure that the cache logic will not cache props info
	instance.Flush()
	return err
}

// ------------------------------------------------------------
// ------------------- db function ----------------------------
// ------------------------------------------------------------

// DBGetOne gets one record from table instance by condition "where"
func DBGetOne(db *sql.DB, instancePtr Model, where map[string]interface{}) (err error) {
	if nil == db {
		return errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildSelect(instancePtr.TableName(), where, nil)
	if nil != err {
		return
	}
	row, err := db.Query(cond, vals...)
	if nil != err || nil == row {
		return
	}
	defer row.Close()
	err = scanner.Scan(row, instancePtr)
	return
}

// DBTxGetOne gets one record from table pre_forum_thread by condition "where" in Tx
func DBTxGetOne(tx *sql.Tx, instancePtr Model, where map[string]interface{}) (err error) {
	if nil == tx {
		return errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildSelect(instancePtr.TableName(), where, nil)
	if nil != err {
		return
	}
	row, err := tx.Query(cond, vals...)
	if nil != err || nil == row {
		return
	}
	defer row.Close()
	err = scanner.Scan(row, instancePtr)
	return
}

// DBGetMulti gets multiple records from table pre_forum_thread by condition "where"
func DBGetMulti(db *sql.DB, instacePtr Model, instaceSlicePtr interface{}, where map[string]interface{}) (err error) {
	if nil == db {
		return errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildSelect(instacePtr.TableName(), where, nil)
	if nil != err {
		return err
	}
	row, err := db.Query(cond, vals...)
	if nil != err || nil == row {
		return err
	}
	defer row.Close()
	err = scanner.Scan(row, instaceSlicePtr)
	return err
}

// DBTxGetMulti gets multiple records from table pre_forum_thread by condition "where" in Tx
func DBTxGetMulti(tx *sql.Tx, instacePtr Model, instaceSlicePtr interface{}, where map[string]interface{}) (err error) {
	if nil == tx {
		return errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildSelect(instacePtr.TableName(), where, nil)
	if nil != err {
		return err
	}
	row, err := tx.Query(cond, vals...)
	if nil != err || nil == row {
		return err
	}
	defer row.Close()
	err = scanner.Scan(row, instaceSlicePtr)
	return
}

// TODO: 返回ID有问题
// DBInsert inserts an array of data into table
func DBInsert(db *sql.DB, instancePtr Model, data []map[string]interface{}) (id int64, err error) {
	if nil == db {
		return 0, errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildInsert(instancePtr.TableName(), data)
	if nil != err {
		return
	}
	result, err := db.Exec(cond, vals...)
	if nil != err || nil == result {
		return
	}
	return result.LastInsertId()
}

// DBTxInsert inserts an array of data into table in Tx
func DBTxInsert(tx *sql.Tx, instancePtr Model, data []map[string]interface{}) (id int64, err error) {
	if nil == tx {
		return 0, errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildInsert(instancePtr.TableName(), data)
	if nil != err {
		return
	}
	result, err := tx.Exec(cond, vals...)
	if nil != err || nil == result {
		return
	}
	return result.LastInsertId()
}

// DBUpdate updates the table
func DBUpdate(db *sql.DB, instancePtr Model, where, data map[string]interface{}) (aff int64, err error) {
	if nil == db {
		return 0, errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildUpdate(instancePtr.TableName(), where, data)
	if nil != err {
		return
	}
	result, err := db.Exec(cond, vals...)
	if nil != err {
		return
	}
	aff, err = result.RowsAffected()
	if aff > 0 {
		instancePtr.Flush()
	}
	return
}

// DBTxUpdate updates the table in Tx
func DBTxUpdate(tx *sql.Tx, instancePtr Model, where, data map[string]interface{}) (aff int64, err error) {
	if nil == tx {
		return 0, errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildUpdate(instancePtr.TableName(), where, data)
	if nil != err {
		return
	}
	result, err := tx.Exec(cond, vals...)
	if nil != err {
		return
	}
	aff, err = result.RowsAffected()
	if aff > 0 {
		instancePtr.Flush()
	}
	return
}

// DBDelete deletes matched records in instance
func DBDelete(db *sql.DB, instancePtr Model, where map[string]interface{}) (aff int64, err error) {
	if nil == db {
		return 0, errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildDelete(instancePtr.TableName(), where)
	if nil != err {
		return
	}
	result, err := db.Exec(cond, vals...)
	if nil != err {
		return
	}
	aff, err = result.RowsAffected()
	if aff > 0 {
		instancePtr.Flush()
	}
	return
}

// DBTxDelete deletes matched records in instance in Tx
func DBTxDelete(tx *sql.Tx, instancePtr Model, where map[string]interface{}) (aff int64, err error) {
	if nil == tx {
		return 0, errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildDelete(instancePtr.TableName(), where)
	if nil != err {
		return
	}
	result, err := tx.Exec(cond, vals...)
	if nil != err {
		return
	}
	aff, err = result.RowsAffected()
	if aff > 0 {
		instancePtr.Flush()
	}
	return
}

// ------------------------------------------------------------
// ------------------- db function ----------------------------
// ------------------------------------------------------------
