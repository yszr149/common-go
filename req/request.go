package req

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func Get() {
	// 和服务端建立连接
	url := "wss://wsaws.okx.com:8443/ws/v5/public"
	log.Printf("connecting to %s", url)
	c, res, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	log.Printf("响应:%s", fmt.Sprint(res))
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte("{\"op\": \"subscribe\", \"args\": [{\"instId\": \"BTC-USDT-SWAP\", \"channel\": \"bbo-tbt\"}]}"))
	if err != nil {
		fmt.Println(err)
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatal(err)
			break
		}
		var msg map[string]interface{}
		err1 := json.Unmarshal([]byte(message), &msg)
		if err != nil {
			log.Fatal(err1)
		}
		value, ok1 := msg["data"]
		if ok1 {
			if reflect.TypeOf(value).Kind() == reflect.Slice {
				s := reflect.ValueOf(value)
				for i := 0; i < s.Len(); i++ {
					ele := s.Index(i)
					tmp := ele.Interface().(map[string]interface{})
					ts := reflect.ValueOf(tmp["ts"]).String()
					// its, _ := strconv.Atoi(ts)
					its, _ := strconv.ParseInt(ts, 10, 64)
					cts := time.Now().UnixNano() / 1e6
					if cts-its > 50 {
						log.Printf("tmp: %d %d %d", cts, its, cts-its)
					}
				}
			}

		}

	}

}
