# 排行榜

#### 設計想法

```
1. 將服務拆成多層，應用DI方式將每層需要的服務注入
2. 主要可分為以下 (詳情可參見下方表格)
    a. usecase
    b. service
    c. repository
    d. domain
    e. dto
    f. cmd
```

| 分層名稱   | 用途                                                  |
| ---------- | ----------------------------------------------------- |
| usecase    | API 入口,在此呼叫 service 並組裝 API 需要的資訊並回傳 |
| service    | 產出固定的資料格式給予 usecase 呼叫並組裝資訊         |
| repository | 主要與資料庫做互動                                    |
| domain     | 每個服務的共用資訊,提供給各層取用                     |
| dto        | 服務的 input output 結構定義                          |
| cmd        | 服務的入口,如有副程式可在這邊新增入口                 |

#### 使用套件

```
1. 排序部份使用redis的zset自動排序
2. 測試部份使用 miniredis的mock redis做測試以防止操作到真正的資料
3. server使用gin
```

#### 服務啟動方式

```
服務啟動採用docker-compose,內含
    a. dispatcher (簡易版,單純模擬multiple server的情況)
    b. server1,server2 ( 服務內容皆相同 )
    c. redie,redis-adminer
```

```
* 請於本資料夾下執行docker-compose指令
    啟動: docker-compose up -d
    查看: docker-compose ps
    log  : docker-compose logs -f $container_name

* 確認 port: 6837,8000 尚未被使用
```

#### 測試方式

```
cmd: go test -v -failfast usecase/board/*
```

#### API 列表

```
domain:
    本機: http://127.0.0.1:8000
```

##### 取得排行榜

```
Method: GET
Path: /api/v1/leaderboard

request:
{
    "top": 10 // 如未帶入top則預設為10
}

response:
{
    "topPlayers": [
        {
            "clientId": "Player2",
            "score": 12
        },
        {
            "clientId": "Player1",
            "score": 11
        }
    ]
}
```

##### 新增分數

```
Method: POST
Path: /api/v1/score

Header:
{
    "ClientId":"player" // 必填
}

＊如帶入相同名稱則覆蓋舊分數

request:
{
    "score": 10 // min=0,max=100
}

response:
{
    "status": "OK"
}
```
