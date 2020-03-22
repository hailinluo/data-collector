package fund

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/hailinluo/data-collector/logger"
	"github.com/hailinluo/data-collector/storage"
	"github.com/hailinluo/data-collector/storage/structs"
	"github.com/hailinluo/data-collector/utils"
	"strconv"
	"strings"
)

type fundCollector struct {
	spec string
}

func NewFundCollector(opts ...Option) *fundCollector {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}

	return &fundCollector{
		spec: options.spec,
	}
}

func (cb *fundCollector) Spec() string {
	return cb.spec
}

func (cb *fundCollector) Run() {
	// 从数据库获取基金公司列表
	companies, err := structs.GetCompanyList(storage.DbEngine())
	if err != nil {
		logger.Errorf("get company list failed. err: %s", err)
		return
	}

	for _, company := range companies {
		listDoc, err := utils.GetDocument(company.FundListUrl)
		if err != nil {
			logger.Errorf("get html document failed. err: %s, url: %v", err, company.FundListUrl)
			continue
		}

		// 开放式基金 & 场内基金
		var openFunds []*structs.Fund
		listDoc.Find("body").
			Find("div.outer_all").
			Find("div.main-content").
			Find("div.first-block").
			Find("div.common-block-con").
			Find("table.ttjj-table").
			Find("tbody").
			Find("tr").
			Each(func(i int, s *goquery.Selection) {
				var fund structs.Fund

				fund.FundCompanyID = company.CompanyID
				fund.FundCompany = company.CompanyName

				// 详情页URL
				fund.DetailPageUrl = s.Find("td.fund-name-code").Find("a.name").AttrOr("href", "")
				if fund.DetailPageUrl == "" {
					logger.Errorf("parsing fund href failed. text: %s", s.Text())
					return
				}

				// 基金名称
				fund.FundName = s.Find("td.fund-name-code").Find("a.name").Text()
				// 基金代码
				fund.FundId = s.Find("td.fund-name-code").Find("a.code").Text()

				openFunds = append(openFunds, &fund)
			})

		// 货币/理财型基金
		var currencyFunds []*structs.Fund
		listDoc.Find("body").
			Find("div.outer_all").
			Find("div.main-content").
			Find("div.third-block").
			Find("div.common-block-con").
			Find("table.ttjj-table").
			Find("tbody").
			Find("tr").
			Each(func(i int, s *goquery.Selection) {
				var fund structs.Fund
				fund.FundCompanyID = company.CompanyID
				fund.FundCompany = company.CompanyName

				// 详情页URL
				fund.DetailPageUrl = s.Find("td.fund-name-code").Find("a.name").AttrOr("href", "")
				if fund.DetailPageUrl == "" {
					logger.Errorf("parsing fund href failed. text: %s", s.Text())
					return
				}

				// 基金名称
				fund.FundName = s.Find("td.fund-name-code").Find("a.name").Text()
				// 基金代码
				fund.FundId = s.Find("td.fund-name-code").Find("a.code").Text()

				currencyFunds = append(currencyFunds, &fund)
			})

		// 填充详细内容
		funds := openFunds
		funds = append(funds, currencyFunds...)
		for _, item := range funds {
			detailDoc, err := utils.GetDocument(item.DetailPageUrl)
			if err != nil {
				logger.Errorf("get html document failed. err: %s, url: %v", err, item.DetailPageUrl)
				continue
			}

			detailDoc.Find("body").
				Find("div.body#body").
				Find("div.wrapper").
				Find("div.wrapper_min").
				Find("div.merchandiseDetail").
				Find("div.fundDetail-main").
				Find("div.fundInfoItem").
				Each(func(i int, s *goquery.Selection) {
					s.Find("div.infoOfFund").
						Find("table").
						Find("tbody").
						Find("tr").
						Find("td").
						Each(func(i int, selection *goquery.Selection) {
							info := selection.Text()
							if strings.Contains(info, "基金类型：") {
								txt := strings.TrimPrefix(info, "基金类型：")
								strs := strings.Split(txt, "|")
								if len(strs) <= 0 {
									logger.Errorf("parsing fund type failed. fundId: %s, info: %s", item.FundId, info)
									return
								}
								item.FundType = strings.TrimSpace(strs[0])
								return
							}
							if strings.Contains(info, "基金规模：") {
								txt := strings.TrimPrefix(info, "基金规模：")
								strs := strings.Split(txt, "亿元")
								if len(strs) <= 0 {
									logger.Errorf("parsing fund scale failed. fundId: %s, info: %s", item.FundId, info)
									return
								}

								item.FundScale, _ = strconv.ParseFloat(strs[0], 10)
								return
							}
							if strings.Contains(info, "基金经理：") {
								item.FundManager = strings.TrimPrefix(info, "基金经理：")
								item.ManagerUrl = selection.Find("a").AttrOr("href", "")
								return
							}
							if strings.Contains(info, "成 立 日：") {
								item.EstablishDate = strings.TrimPrefix(info, "成 立 日：")
								return
							}
						})

					s.Find("div.infoOfFund").
						Find("table").
						Find("tbody").
						Find("tr").
						Find("td.specialData").
						Each(func(i int, selection *goquery.Selection) {
							info := selection.Text()

							strs := strings.Split(info, "跟踪误差：")
							if len(strs) != 2 {
								logger.Errorf("parsing track failed, info: %s", info)
								return
							}

							// 跟踪误差
							txt := strings.TrimSuffix(strs[1], "%")
							item.TrackDeviation, _ = strconv.ParseFloat(txt, 10)

							// 跟踪标的
							item.TrackTarget = strings.TrimPrefix(strs[0], "跟踪标的：")
							item.TrackTarget = strings.TrimSuffix(item.TrackTarget, " ")
							item.TrackTarget = strings.TrimSuffix(item.TrackTarget, "|")
							item.TrackTarget = strings.TrimSuffix(item.TrackTarget, " ")
						})

					s.Find("div.dataOfFund").
						Find("dl").
						Find("dd").
						Each(func(i int, selection *goquery.Selection) {
							info := selection.Text()
							if strings.Contains(info, "近1月：") {
								yield := strings.TrimPrefix(info, "近1月：")
								yield = strings.TrimSuffix(yield, "%")
								item.YieldMonth, _ = strconv.ParseFloat(yield, 10)
							}
							if strings.Contains(info, "近3月：") {
								yield := strings.TrimPrefix(info, "近3月：")
								yield = strings.TrimSuffix(yield, "%")
								item.Yield3Month, _ = strconv.ParseFloat(yield, 10)
							}
							if strings.Contains(info, "近6月：") {
								yield := strings.TrimPrefix(info, "近6月：")
								yield = strings.TrimSuffix(yield, "%")
								item.Yield6Month, _ = strconv.ParseFloat(yield, 10)
							}
							if strings.Contains(info, "近1年：") {
								yield := strings.TrimPrefix(info, "近1年：")
								yield = strings.TrimSuffix(yield, "%")
								item.YieldYear, _ = strconv.ParseFloat(yield, 10)
							}
							if strings.Contains(info, "近3年：") {
								yield := strings.TrimPrefix(info, "近3年：")
								yield = strings.TrimSuffix(yield, "%")
								item.Yield3Year, _ = strconv.ParseFloat(yield, 10)
							}
							if strings.Contains(info, "成立来：") {
								yield := strings.TrimPrefix(info, "成立来：")
								yield = strings.TrimSuffix(yield, "%")
								item.Yield, _ = strconv.ParseFloat(yield, 10)
							}
						})

					nav := s.Find("div.dataOfFund").
						Find("dl.dataItem02").
						Find("dd.dataNums").
						Find("span.ui-font-large").Text()
					item.NAV, _ = strconv.ParseFloat(nav, 10)
					anv := s.Find("div.dataOfFund").
						Find("dl.dataItem03").
						Find("dd.dataNums").
						Find("span.ui-font-large").Text()
					item.ANV, _ = strconv.ParseFloat(anv, 10)
				})
			fmt.Println(item)
		}

		// 插入或更新数据库
		for _, item := range funds {
			exist, err := structs.ExistFund(storage.DbEngine(), item.FundId)
			if err != nil {
				logger.Errorf("check exist failed. err: %s, item: %v", err, item)
				continue
			}

			// update
			if exist {
				err = structs.UpdateFund(storage.DbEngine(), item.FundId, item)
				if err != nil {
					logger.Errorf("update record failed. err: %s, item: %v", err, item)
				}
				continue
			}

			// insert
			err = structs.AddFund(storage.DbEngine(), item)
			if err != nil {
				logger.Errorf("insert failed. err: %s, item: %v", err, item)
			}
		}
		break
	}
}
