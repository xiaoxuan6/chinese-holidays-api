package handlers

import (
    "encoding/json"
    "github.com/bastengao/chinese-holidays-go/holidays"
    "github.com/godcong/chronos"
    "github.com/gorilla/mux"
    "github.com/huandu/xstrings"
    "net/http"
    "strings"
    "time"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write([]byte(`{"code": 500, "msg": "route [` + r.RequestURI + `] not found!"}`))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write([]byte(`{"code": 200, "msg": "中国特色的休假安排或者工作日查询 api 接口，所有数据均来自国务院发布，\n 源码地址：https://github.com/xiaoxuan6/chinese-holidays-api"}`))
}

var (
    zhYear = map[string]string{
        "子年": "鼠年",
        "丑年": "牛年",
        "寅年": "虎年",
        "卯年": "兔年",
        "辰年": "龙年",
        "巳年": "蛇年",
        "午年": "马年",
        "未年": "羊年",
        "申年": "猴年",
        "酉年": "鸡年",
        "戌年": "狗年",
        "亥年": "猪年",
    }

    weekly = []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}
)

func DateHandler(w http.ResponseWriter, r *http.Request) {
    date := r.URL.Query().Get("date")
    if date == "" {
        vars := mux.Vars(r)
        date = vars["date"]
    }

    if len(date) < 1 {
        w.Header().Set("Content-Type", "application/json")
        _, _ = w.Write([]byte(`{"code": 500, "msg": "参数：date 不能为空！"}`))
        return
    }

    parsedDate, err := time.Parse(time.DateOnly, date)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        _, _ = w.Write([]byte(`{"code": 500, "msg": "Invalid date format"}`))
        return
    }

    queryer, _ := holidays.BundleQueryer()
    holiday, err := queryer.IsHoliday(parsedDate)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        _, _ = w.Write([]byte(`{"code": 500, "msg": "query holiday fail: "` + err.Error() + `}`))
        return
    }

    workingday, err := queryer.IsWorkingday(parsedDate)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        _, _ = w.Write([]byte(`{"code": 500, "msg": "query working day fail: "` + err.Error() + `}`))
        return
    }

    lunarDate := chronos.New(parsedDate).LunarDate()
    lunarYear := xstrings.Slice(lunarDate, 0, 3)
    response := map[string]interface{}{
        "code":           200,
        "msg":            "query ok",
        "date":           parsedDate.Format(time.DateOnly),
        "lunar_date":     strings.ReplaceAll(lunarDate, lunarYear, ""),
        "lunar_year":     lunarYear,
        "is_holiday":     holiday,
        "is_working_day": workingday,
        "weekday":        weekly[parsedDate.Weekday()],
        "zh_year":        zhYear[xstrings.Slice(lunarDate, 1, 3)],
    }

    b, _ := json.Marshal(response)

    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write(b)

    return
}
