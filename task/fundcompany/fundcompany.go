package fundcompany

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/hailinluo/data-collector/logger"
	"github.com/hailinluo/data-collector/storage"
	"github.com/hailinluo/data-collector/storage/structs"
	"github.com/hailinluo/data-collector/utils"
	"strconv"
	"strings"
)

type fcCollector struct {
	spec     string
	homePage string
	resURl   string
}

func NewFcCollector(opts ...Option) *fcCollector {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}

	return &fcCollector{
		spec:     options.spec,
		homePage: options.homePage,
		resURl:   options.resUrl,
	}
}

func (cb *fcCollector) Spec() string {
	return cb.spec
}

func (cb *fcCollector) Run() {
	// 获取基金公司列表
	listDoc, err := utils.GetDocument(cb.resURl)
	if err != nil {
		logger.Errorf("get html document failed. err: %s, url: %v", err, cb.resURl)
		return
	}

	var companyList []*structs.CompanyInfo
	listDoc.Find("div.outer_all").
		Find("div.ttjj-grid-row").
		Find("div.main-content").
		Find("div.fourth-block").
		Find("div#companyCon").
		Find("div#companyTable.common-block-con").
		Find("table#gspmTbl.ttjj-table").
		Find("tbody").Find("tr").
		Each(func(i int, s *goquery.Selection) {
			var company structs.CompanyInfo

			companySelection := s.Find("td.td-align-left").Find("a")
			company.PageUrl = companySelection.AttrOr("href", "")
			if company.PageUrl != "" {
				company.CompanyName = companySelection.Text()

				// 解析ID
				strs := strings.Split(company.PageUrl, "/")
				company.CompanyID = strings.Trim(strs[len(strs)-1], ".html")

				company.PageUrl = cb.homePage + company.PageUrl

				companyList = append(companyList, &company)
			}
		})

	for _, item := range companyList {
		itemDoc, err := utils.GetDocument(item.PageUrl)
		if err != nil {
			logger.Errorf("get html document failed. err: %s, url: %v", err, item.PageUrl)
			continue
		}

		itemDoc.Find("body").
			Find("div.outer_all").
			Find("div.ttjj-grid-row").
			Find("div.main-content").
			Find("div.common-basic-info").
			Find("div.fund-info").
			Find("ul").
			Each(func(i int, s *goquery.Selection) {
				// 基金规模
				scale := s.Find("li.padding-left-10").Find("label").Text()
				scale = strings.Trim(scale, "亿元")
				item.AUM, _ = strconv.ParseFloat(scale, 10)

				// 成立日期
				item.EstablishDate = s.Find("li.date").Find("label").Text()

				// 基金列表页URL
				item.FundListUrl = cb.homePage + "/Company/f10/jjjz_" + item.CompanyID + ".html"
		})
	}

	// 插入或更新数据库
	for _, item := range companyList {
		exist, err := structs.ExistAcompanyInfo(storage.DbEngine(), item.CompanyID)
		if err != nil {
			logger.Errorf("check exist failed. err: %s, item: %v", err, item)
			continue
		}

		// update
		if exist {
			err = structs.UpdateCompany(storage.DbEngine(), item.CompanyID, item)
			if err != nil {
				logger.Errorf("update record failed. err: %s, item: %v", err, item)
			}
			continue
		}

		// insert
		err = structs.AddCompany(storage.DbEngine(), item)
		if err != nil {
			logger.Errorf("insert failed. err: %s, item: %v", err, item)
		}
	}
}
