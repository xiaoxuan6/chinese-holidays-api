package main

import (
	"github.com/bastengao/chinese-holidays-go/holidays"
	"github.com/godcong/chronos"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))

	r.NotFoundHandler = http.HandlerFunc(notFoundHandle)

	r.HandleFunc("/", indexHandle).Methods(http.MethodGet)

	r.HandleFunc("/api/holidays", dateHandle).Methods(http.MethodGet).Queries("date", "{date}")
	r.HandleFunc("/api/holidays/{date}", dateHandle).Methods(http.MethodGet)

	log.Panic(http.ListenAndServe(":80", r))
}

func notFoundHandle(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"code": 500, "msg": "route [` + r.RequestURI + `] not found!"}`))
}

func indexHandle(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"code": 200, "msg": "中国特色的休假安排或者工作日查询 api 接口，所有数据均来自国务院发布，\n 源码地址：https://github.com/xiaoxuan6/chinese-holidays-api"}`))
}

func dateHandle(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		vars := mux.Vars(r)
		date = vars["date"]
	}

	if len(date) < 1 {
		r.Header.Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code": 500, "msg": "参数：date 不能为空！"}`))
		return
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		r.Header.Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code": 500, "msg": "Invalid date format"}`))
		return
	}

	queryer, _ := holidays.BundleQueryer()
	holiday, err := queryer.IsHoliday(parsedDate)
	if err != nil {
		r.Header.Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code": 500, "msg": "query holiday fail: "` + err.Error() + `}`))
		return
	}

	workingday, err := queryer.IsWorkingday(parsedDate)
	if err != nil {
		r.Header.Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code": 500, "msg": "query working day fail: "` + err.Error() + `}`))
		return
	}

	lunarDate := chronos.New(parsedDate).LunarDate()
	r.Header.Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"code": 200, "msg": "query ok", "date": "` + parsedDate.Format(time.DateOnly) + `", "lunar_date": "` + lunarDate + `", "is_holiday": ` + strconv.FormatBool(holiday) + `, "is_working_day": ` + strconv.FormatBool(workingday) + `, "weekday": "` + weekdayZh(parsedDate) + `"}`))

	return
}

func weekdayZh(date time.Time) (weekdayStr string) {
	weekday := date.Weekday()
	switch weekday {
	case time.Monday:
		weekdayStr = "星期一"
	case time.Tuesday:
		weekdayStr = "星期二"
	case time.Wednesday:
		weekdayStr = "星期三"
	case time.Thursday:
		weekdayStr = "星期四"
	case time.Friday:
		weekdayStr = "星期五"
	case time.Saturday:
		weekdayStr = "星期六"
	case time.Sunday:
		weekdayStr = "星期日"
	}

	return
}
