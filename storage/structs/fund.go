package structs

import (
	"github.com/go-xorm/xorm"
	"github.com/hailinluo/data-collector/logger"
)

type Fund struct {
	FundId         string  `xorm:"'fund_id' pk" json:"FundId"`
	FundName       string  `xorm:"'fund_name'" json:"FundName"`
	FundScale      float64 `xorm:"'fund_scale' comment('基金规模')" json:"FundScale"`
	FundType       string  `xorm:"'fund_type' comment('基金类型')" json:"FundType"`
	FundManager    string  `xorm:"'fund_manager' comment('基金经理')" json:"FundManager"`
	ManagerUrl     string  `xorm:"'manager_url' comment('基金经理详情页URL')" json:"ManagerUrl"`
	EstablishDate  string  `xorm:"'establish_date' comment('成立日期')" json:"EstablishDate"`
	TrackTarget    string  `xorm:"'track_target' comment('跟踪标的')" json:"TrackTarget"`
	TrackDeviation float64  `xorm:"'track_deviation' comment('跟随误差')" json:"TrackDeviation"`
	DetailPageUrl  string  `xorm:"'detail_page_url' comment('详情页URL')" json:"DetailPageUrl"`
	FundCompanyID  string  `xorm:"'fund_company_id' comment('基金公司编号')" json:"FundCompanyID"`
	FundCompany    string  `xorm:"'fund_company' comment('基金公司')" json:"FundCompany"`
	NAV            float64  `xorm:"'nav' comment('单位净值')" json:"NAV"`
	ANV            float64  `xorm:"'anv' comment('累计净值')" json:"ANV"`
	YieldMonth     float64 `xorm:"'yield_month' comment('近1月收益率')" json:"YieldMonth"`
	Yield3Month    float64 `xorm:"'yield_3month' comment('近3月收益率')" json:"Yield3Month"`
	Yield6Month    float64 `xorm:"'yield_6month' comment('近半年收益率')" json:"Yield6Month"`
	YieldYear      float64 `xorm:"'yield_year' comment('近1年收益率')" json:"YieldYear"`
	Yield3Year     float64 `xorm:"'yield_3year' comment('近3年收益率')" json:"Yield3Year"`
	Yield          float64 `xorm:"'yield' comment('成立以来收益率')" json:"Yield"`
}

func (*Fund) TableName() string {
	return "t_fund"
}

func ExistFund(engine *xorm.Engine, id string) (bool, error) {
	exist, err := engine.Exist(&Fund{FundId: id})
	if err != nil {
		logger.Errorf("check exist failed. err: %s", err)
		return false, nil
	}
	return exist, nil
}

func AddFund(engine *xorm.Engine, fund *Fund) error {
	_, err := engine.InsertOne(fund)
	if err != nil {
		return err
	}
	return nil
}

func AddFundList(engine *xorm.Engine, funds []*Fund) error {
	session := engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	inserted, err := session.Insert(funds)
	if err != nil {
		return err
	}

	if inserted != int64(len(funds)) {
		return session.Rollback()
	}

	return session.Commit()
}

func UpdateFund(engine *xorm.Engine, fundId string, fund *Fund) error {
	_, err := engine.Update(fund, &Fund{FundId: fundId})
	if err != nil {
		return err
	}
	return nil
}

func QueryFund(engine *xorm.Engine, fundId string) (*Fund, error) {
	var ret Fund
	err := engine.Where("fund_id = ?", fundId).Find(&ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetFundList(engine *xorm.Engine) ([]*Fund, error) {
	var companies []*Fund
	err := engine.Find(&companies)
	if err != nil {
		return nil, err
	}
	return companies, nil
}
