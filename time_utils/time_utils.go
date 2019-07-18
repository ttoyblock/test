package time_utils

import (
	"time"
)

var BeijingLocation = time.FixedZone("Asia/Shanghai", 8*60*60)

const (
	RFC3339 = "2006-01-02T15:04:05+08:00"
	TT      = "2006-01-02 15:04:05"
)

// GetNow 获取当前时间 datetime
func GetNow() time.Time {
	return time.Now().In(BeijingLocation)
}

// GetNowDateTime 获取当前时间 2019-07-12 20:51:36
func GetNowDateTime() string {
	return time.Now().In(BeijingLocation).Format(TT)
}

// GetBeforeTimeStr 获取beforeSecond前的时间
// return 2019-07-12 20:52:26
func GetBeforeTimeStr(beforeSecond int64) string {
	return time.Unix(time.Now().In(BeijingLocation).Unix()-beforeSecond, 0).Format(TT)
}

// GetNextTimeStr 获取nextSecond之后的时间
// return 2019-07-12 20:52:26
func GetNextTimeStr(nextSecond int) string {
	return time.Unix(time.Now().In(BeijingLocation).Unix()+int64(nextSecond), 0).Format(TT)
}

// GetNextTime 2019-07-12 20:56:01 +0800 CST
func GetNextTime(nextSecond int) time.Time {
	return time.Unix(time.Now().In(BeijingLocation).Unix()+int64(nextSecond), 0)
}

// GetNowDate 2006-01-02
func GetNowDate() string {
	return time.Now().In(BeijingLocation).Format("2006-01-02")
}

// GetNowDayNum 190712
func GetNowDayNum() string {
	return time.Now().In(BeijingLocation).Format("060102")
}

// GetNowUnix 获取当前时间戳
func GetNowUnix() int64 {
	return time.Now().In(BeijingLocation).Unix()
}

// GetDayStart 获取零晨时间
// return 2019-07-12 00:00:00
func GetDayStart(t time.Time) string {
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, BeijingLocation)
	return tm1.Format(TT)
}

// GetDayStartTime 2019-07-14 00:00:00 +0800 Asia/Shanghai
func GetDayStartTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, BeijingLocation)
}

// GetDayEnd 获取一个最后一刻时间
// return 2019-07-14 23:59:59
func GetDayEnd(t time.Time) string {
	tm := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, BeijingLocation)
	return tm.Format(TT)
}

// GetNowDateTimeFmt 获取当前时间
// return 20190714101002
func GetNowDateTimeFmt() string {
	return time.Now().In(BeijingLocation).Format("20060102150405")
}

// GetNowDayEndTS 获取当天最后时间时间戳
func GetNowDayEndTS() int64 {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, BeijingLocation).Unix()
}

// GetYesterdayEndTime 获取昨天最后时间
// return 2019-07-13 23:59:59
func GetYesterdayEndTime() string {
	t := time.Unix(time.Now().In(BeijingLocation).Unix()-int64(86400), 0)
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, BeijingLocation).Format(TT)
}

// GetYesterdayStartTime 获取昨天开始时间
// return 2019-07-13 00:00:00
func GetYesterdayStartTime() string {
	t := time.Unix(time.Now().In(BeijingLocation).Unix()-int64(86400), 0)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, BeijingLocation).Format(TT)
}

// DataTimeStr2TS 2006-01-02 15:04:05 转时间戳
func DataTimeStr2TS(str string) int64 {
	tm, err := time.ParseInLocation(TT, str, BeijingLocation)
	if nil != err {
		return int64(0)
	}
	return tm.Unix()
}

// RFC3339StrT2TS RFC3339 to ts
// param 2006-01-02T15:04:05+08:00
func RFC3339StrT2TS(layout string) int64 {
	theTime, _ := time.ParseInLocation(time.RFC3339, layout, BeijingLocation)
	return theTime.Unix()
}

// ParseDateTime 将2006-01-02 15:04:05转成time.Time
func ParseDateTime(layout string) (time.Time, error) {
	t, err := time.ParseInLocation(TT, layout, BeijingLocation)
	if err != nil {
		return t, err
	}
	return t, nil
}

// DateTime2DateTimeStr _
// return 2006-01-02 15:04:00
func DateTime2DateTimeStr(t time.Time) string {
	return t.In(BeijingLocation).Format(TT)
}

// DateTimeToTime _
// param 2006-01-02 15:04:05
// return 2006-01-02 15:04:05 +0800 Asia/Shanghai
func DateTimeToTime(str string) (time.Time, error) {
	tm, err := time.ParseInLocation(TT, str, BeijingLocation)
	if nil != err {
		return tm, err
	}
	return tm, nil
}

// RFC3339ToTT RFC格式转成正常格式
// param 2006-01-02T15:04:05+08:00
// return 2006-01-02 15:04:05
func RFC3339ToTT(s string) string {
	t, err := time.ParseInLocation(RFC3339, s, BeijingLocation)
	if err != nil {
		return "0000-00-00 00:00:00"
	}
	return t.Format(TT)
}

// DateTime2DateStr 获取时间日期
// return 2019-07-14
func DateTime2DateStr(t time.Time) string {
	return t.In(BeijingLocation).Format("2006-01-02")
}

// TS2DateTime 时间戳转时间
// return 2019-07-14 10:21:59
func TS2DateTime(second int64) string {
	return time.Unix(second, 0).In(BeijingLocation).Format(TT)
}
