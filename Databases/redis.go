package Databases

import (
	"github.com/gomodule/redigo/redis"
)

var RedisPool *redis.Pool
var RedisConn redis.Conn

func RedisPollInit() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   3,
		MaxActive: 5,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}

func RedisInit() {
	RedisPool = RedisPollInit()
	RedisConn = RedisPool.Get()
}

func RedisClose() {
	_ = RedisPool.Close()
}

//在redis操作golang 结构体切片示例
//func HandleStructSlice() {
//	P := Models.Product{Name: sql.NullString{String: "草", Valid: true}}
//	moduleList := []Models.Product{P, P, P, P}
//	datas, _ := json.Marshal(moduleList)
//
//	_, _ = RedisConn.Do("SET", "ModuleList", datas)
//	rebytes, _ := redis.Bytes(RedisConn.Do("get", "ModuleList"))
//	var object []*Models.Product
//	_ = json.Unmarshal(rebytes, &object)
//	for _, v := range object {
//		fmt.Printf("%+v\n", *v)
//	}
//}
//
//在redis操作golang结构体示例
//func HandleStruct() {
//	var testStruct = &Models.Product{Name: sql.NullString{String: "cx", Valid: true}}
//	//json序列化
//	datas, _ := json.Marshal(testStruct)
//	//缓存数据
//	_, _ = RedisConn.Do("set", "struct3", datas)
//	//读取数据
//	rebytes, _ := redis.Bytes(RedisConn.Do("get", "struct3"))
//	//json反序列化
//	object := &Models.Product{}
//	_ = json.Unmarshal(rebytes, object)
//	fmt.Println(object)
//}
