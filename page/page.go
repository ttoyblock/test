package page

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func call(param url.Values, res interface{}) error {
	req, _ := http.NewRequest(http.MethodPost, "https://www.lagou.com/jobs/positionAjax.json?needAddtionalResult=false&isSchoolJob=0", strings.NewReader(param.Encode()))
	req.Header.Set("Origin", "https://www.lagou.com")
	req.Header.Set("X-Anit-Forge-Code", "0")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Referer", "https://www.lagou.com/jobs/list_go?city=%E5%85%A8%E5%9B%BD&cl=false&fromSearch=true&labelWords=&suginput=")
	req.Header.Set("X-Anit-Forge-Token", "None")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if res != nil {
		d := json.NewDecoder(resp.Body)
		d.Decode(res)
	}

	return nil
}

// 由于 拉勾 只有30页。所以，这个地方就直接写死 30 了
const pagesize = 30

func Lagou(language string) (res [][]string) {
	param := url.Values{}
	param.Set("first", "true")
	param.Set("kd", language)

	workInfos := make([]WorkInfo, 0, 400)

	for i := 1; i <= pagesize; i++ {
		fmt.Printf("now is %d page \n", i)
		param.Set("pn", strconv.Itoa(i))
		ret := LagouRet{}
		err := call(param, &ret)
		if err != nil {
			fmt.Println("err: ", err)
			return
		}

		if ret.Code != 0 {
			err = fmt.Errorf("%s", ret.Msg)
			return
		}

		workInfos = append(workInfos, ret.Content.PositionResult.Result...)
		param.Del("first")
		param.Set("first", "false")
		time.Sleep(10 * time.Second)
	}

	return ConvertWorkInfo2Csv(workInfos)
}

func ConvertWorkInfo2Csv(infos []WorkInfo) (res [][]string) {
	title := []string{"工作年限", "教育水平", "工资", "职位名称", "职位标签", "地点", "公司名称", "公司类型", "融资情况", "公司标签"}
	res = make([][]string, 0)
	res = append(res, title)
	tmp := make([]string, 10)
	for _, v := range infos {
		tmp = []string{v.WorkYear, v.Education, v.Salary, v.PositionName, strings.Join(v.PositionLables, ":"), v.City + " " +
			strings.Join(v.BusinessZones, ":"), v.CompanyFullName, v.IndustryField, v.FinanceStage, strings.Join(v.CompanyLabelList, ":")}
		res = append(res, tmp)
	}

	return
}

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type LagouRet struct {
	Error
	Content struct {
		HrMap          map[string]HrInfo `json:"hrInfoMap"`
		PositionResult struct {
			Result []WorkInfo `json:"result"`
		} `json:"positionResult"`
	} `json:"content"`
}

type HrInfo struct {
	CanTalk      bool   `json:"canTalk"`
	Phone        string `json:"phone"`
	Portrait     string `json:"portrait"`
	PositionName string `json:"positionName"`
	RealName     string `json:"realName"`
	ReceiveEmail string `json:"receiveEmail"`
	UserLevel    string `json:"userLevel"`
}

type WorkInfo struct {
	WorkYear              string      `json:"workYear"`
	Education             string      `json:"education"`
	JobNature             string      `json:"jobNature"`
	CompanyID             int         `json:"companyId"`
	PositionName          string      `json:"positionName"`
	PositionID            int         `json:"positionId"`
	CreateTime            string      `json:"createTime"`
	City                  string      `json:"city"`
	CompanyLogo           string      `json:"companyLogo"`
	IndustryField         string      `json:"industryField"`
	PositionAdvantage     string      `json:"positionAdvantage"`
	Salary                string      `json:"salary"`
	CompanySize           string      `json:"companySize"`
	Approve               int         `json:"approve"`
	Score                 int         `json:"score"`
	CompanyShortName      string      `json:"companyShortName"`
	PositionLables        []string    `json:"positionLables"`
	IndustryLables        []string    `json:"industryLables"`
	PublisherID           int         `json:"publisherId"`
	FinanceStage          string      `json:"financeStage"`
	CompanyLabelList      []string    `json:"companyLabelList"`
	District              string      `json:"district"`
	BusinessZones         []string    `json:"businessZones"`
	FormatCreateTime      string      `json:"formatCreateTime"`
	AdWord                int         `json:"adWord"`
	CompanyFullName       string      `json:"companyFullName"`
	ImState               string      `json:"imState"`
	LastLogin             int64       `json:"lastLogin"`
	Explain               interface{} `json:"explain"`
	Plus                  interface{} `json:"plus"`
	PcShow                int         `json:"pcShow"`
	AppShow               int         `json:"appShow"`
	Deliver               int         `json:"deliver"`
	GradeDescription      interface{} `json:"gradeDescription"`
	PromotionScoreExplain interface{} `json:"promotionScoreExplain"`
	FirstType             string      `json:"firstType"`
	SecondType            string      `json:"secondType"`
	IsSchoolJob           int         `json:"isSchoolJob"`
}

func writeCsv2File(result [][]string, fileName string) {
	csvBuf := new(bytes.Buffer)
	writter := csv.NewWriter(csvBuf)

	writter.WriteAll(result)
	writter.Flush()

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()
	f.Write(csvBuf.Bytes())
}

func main() {
	// 这个地方的参数是你想要查的职位相关的东西
	writeCsv2File(Lagou("go"), "lagou.csv")
}
