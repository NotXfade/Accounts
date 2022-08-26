package accounts

import (
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/util"
)

//batch is used for getting data
type batch struct {
	CreatedAt time.Time
}

type batchDate struct {
	CreatedAt string
}

//BatchList is used to get batch list
func BatchList() (map[string]interface{}, int) {
	mapd := make(map[string]interface{})
	db := config.DB
	//db = db.Debug()
	list := []batch{}
	query := "select distinct cast(created_at as date) from accounts"
	db.Raw(query).Scan(&list)
	//layoutUS := "January 2, 2006"
	layoutISO := "2006-01-02"
	batchlist := []batchDate{}
	for i := 0; i < len(list); i++ {
		//log.Println(list[i].CreatedAt)
		batchlist = append(batchlist, batchDate{
			CreatedAt: list[i].CreatedAt.Format(layoutISO),
		})
	}
	mapd["error"] = false
	mapd["message"] = "Operation successfull"
	mapd["list"] = batchlist
	return mapd, 200
}

//ReportListing is used to
type ReportListing struct {
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Marks float64 `json:"marks"`
	Link  string  `json:"link"`
}

//GetReports is used to get reports listing
func GetReports(level, batch, module, sortby, limit, page string) (map[string]interface{}, int) {
	mapd := make(map[string]interface{})
	db := config.DB
	reports := []ReportListing{}
	limits, _ := strconv.Atoi(limit)
	pageno, _ := strconv.Atoi(page)
	offset := (limits * (pageno - 1))
	accounts := []database.Accounts{}
	var count int
	if batch == "" {
		db.Model(&database.Accounts{}).Where("level=? AND account_status=?", level, "active").Count(&count)
	} else {
		type countList struct {
			Count int
		}
		countlist := countList{}
		queryString := "select count(*) from accounts where level='" + level + "' AND account_status='active' " + "AND  created_at ::date ='" + batch + "' "
		db.Raw(queryString).Scan(&countlist)
		count = countlist.Count
	}
	query := "select * from accounts where level='" + level + "' AND account_status='active' "
	if batch != "" {
		query = query + "AND  created_at ::date ='" + batch + "' "
	}
	query = query + "ORDER BY lower(name) " + "LIMIT " + limit + " OFFSET " + strconv.Itoa(offset) + ";"
	db.Raw(query).Scan(&accounts)
	for i := 0; i < len(accounts); i++ {
		personal := []database.PersonalInfo{}
		marks := RequestReports(accounts[i].Email, module)
		db.Where("userid=?", accounts[i].ID).Find(&personal)
		var link string
		var err error
		if len(personal) != 0 {
			if personal[0].ProfileImgLink != "" {
				link, err = SendRequestForLink(personal[0].ProfileImgLink)
				if err != nil {
					log.Println(err)
				}
			}
		}
		if marks != -1 {
			reports = append(reports, ReportListing{
				Name:  accounts[i].Name,
				Email: accounts[i].Email,
				Marks: marks,
				Link:  link,
			})
		}
	}
	if sortby == "scoring" {
		sort.Slice(reports, func(i, j int) bool {
			return reports[i].Marks > reports[j].Marks
		})
	}
	mapd["error"] = false
	mapd["message"] = "Operation successfull"
	mapd["list"] = reports
	mapd["count"] = count
	return mapd, 200
}

type payload struct {
	Email  string `json:"email"`
	Module string `json:"module"`
}

//RequestReports is
func RequestReports(email, module string) (marks float64) {
	defer util.Panic()
	data := payload{
		Email:  email,
		Module: module,
	}
	sendData, err := json.Marshal(data)
	log.Println(err)
	nc := config.NC
	res, err := nc.Request(config.Conf.NatsServer.Subject+".module.scores", sendData, 1*time.Minute)
	if err != nil {
		log.Println(err)
		return 0
	}
	type response struct {
		Marks float64 `json:"scores"`
	}
	resp := response{}
	err = json.Unmarshal(res.Data, &resp)
	if err != nil {
		log.Println(err)
		return 0
	}
	return resp.Marks
}
