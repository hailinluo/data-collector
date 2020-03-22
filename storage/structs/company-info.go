package structs

import (
	"github.com/go-xorm/xorm"
	"github.com/hailinluo/data-collector/logger"
	"time"
)

type CompanyInfo struct {
	CompanyID     string  `xorm:"'company_id' pk" json:"CompanyID"`
	CompanyName   string  `xorm:"'company_name'" json:"CompanyName"`
	AUM           float64 `xorm:"'aum' comment('资产管理规模')" json:"AUM"`
	EstablishDate string  `xorm:"'establish_date' comment('成立日期')" json:"EstablishDate"`
	PageUrl       string  `xorm:"'page_url' comment('天天基金网页面地址')" json:"PageUrl"`
	FundListUrl   string  `xorm:"'fund_list_url' comment('基金列表页面地址')" json:"FundListUrl"`
	// TODO active

	Created time.Time `xorm:"created" json:"Created"`
	Updated time.Time `xorm:"updated" json:"Updated"`
}

func (*CompanyInfo) TableName() string {
	return "t_company_info"
}

func ExistAcompanyInfo(engine *xorm.Engine, id string) (bool, error) {
	exist, err := engine.Exist(&CompanyInfo{CompanyID:id})
	if err != nil {
		logger.Errorf("check exist failed. err: %s", err)
		return false, nil
	}
	return exist, nil
}

func AddCompany(engine *xorm.Engine, company *CompanyInfo) error {
	_, err := engine.InsertOne(company)
	if err != nil {
		return err
	}
	return nil
}

func AddCompanyList(engine *xorm.Engine, companies []*CompanyInfo) error {
	session := engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	inserted, err := session.Insert(companies)
	if err != nil {
		return err
	}

	if inserted != int64(len(companies)) {
		return session.Rollback()
	}

	return session.Commit()
}

func UpdateCompany(engine *xorm.Engine, companyId string, company *CompanyInfo) error {
	_, err := engine.Update(company, &CompanyInfo{CompanyID: companyId})
	if err != nil {
		return err
	}
	return nil
}

func QueryCompany(engine *xorm.Engine, companyId string) (*CompanyInfo, error) {
	var ret CompanyInfo
	err := engine.Where("company_id = ?", companyId).Find(&ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetCompanyList(engine *xorm.Engine) ([]*CompanyInfo, error) {
	var companies []*CompanyInfo
	err := engine.Desc("aum").Find(&companies)
	if err != nil {
		return nil, err
	}
	return companies, nil
}
