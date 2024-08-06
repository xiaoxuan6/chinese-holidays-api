# chinese-holidays-api

提供中国休假或者工作日查询

## Docker

```docker
docker run --name=holidays -p 80:80 -d ghcr.io/xiaoxuan6/chinese-holidays-api/chinese-holidays-api:latest
```

## Api

查询 `2024-08-06` 是否是工作日

```shell
curl http://127.0.0.1/api/holidays/2024-08-06
or
curl http://127.0.0.1/api/holidays?date=2024-08-06

{
  "code": 200,
  "msg": "query ok",
  "is_holiday": false,
  "is_working_day": true
}
```

返回参数：

|字段|描述|
|:---|:---|
|code|返回状态值（200：表示成功，500：表示失败）|
|msg|描述|
|is_holiday|是否是节假日|
|is_working_day|是否是工作日|

## Features

- [x] bundled data
    - support [2024](https://www.gov.cn/zhengce/content/202310/content_6911527.htm)
    - support [2023](http://www.gov.cn/zhengce/content/2022-12/08/content_5730844.htm)
    - support [2022](http://www.gov.cn/zhengce/content/2021-10/25/content_5644835.htm)
    - support [2021](http://www.gov.cn/zhengce/content/2020-11/25/content_5564127.htm)
    - support [2020](http://www.gov.cn/zhengce/content/2019-11/21/content_5454164.htm)
    - support [2019](http://www.gov.cn/zhengce/content/2018-12/06/content_5346276.htm) and
      5.1 [changes](http://www.gov.cn/zhengce/content/2019-03/22/content_5375877.htm)
    - support [2018](http://www.gov.cn/zhengce/content/2017-11/30/content_5243579.htm)
    - support [2017](http://www.gov.cn/zhengce/content/2016-12/01/content_5141603.htm)
    - support 2016
