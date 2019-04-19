package nmeasentencetomqtt_flogo

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/go-redis/redis"
	"strings"
	"fmt"
	"strconv"
	"encoding/json"
	"math"
)

var log = logger.GetLogger("nmeasentencetomqtt_log")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}
type Variable struct {
	Key string `json:"key"`
	Desc  string `json:"desc"`
	Indexes []struct {
				Index int `json:"index"`
				Name string `json:"name"`
				Type string `json:"type"`
		} `json:"indexes"`
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error)  {
	address := context.GetInput("address").(string)
	dbNo := context.GetInput("dbNo").(int)
	sentence := context.GetInput("sentence").(string)
	log.Infof("Connecting to Redis: [%s]", address)
	log.Infof("DB no: [%s]",dbNo)

	splitted := strings.Split(sentence, ",")

	client := redis.NewClient(&redis.Options{
			Addr:     address,
			Password: "", // no password set
			DB:       dbNo,  // use default DB
		})

		val, err := client.Get(splitted[0]).Result()
		mqtt_message := ""
		topic := ""
		m := make(map[string]interface{})

		if err != nil {
			log.Errorf("Error! Not found Redis var: %s", splitted[0] )
		} else {
			bytes := []byte(val)

			var v Variable
				err2 := json.Unmarshal(bytes, &v)
				if err2 != nil {
					panic(err2)
				}
			topic = v.Key
			m["desc"] = v.Desc
			 for _, val := range v.Indexes {
				 if val.Type == "string" {
					 m[val.Name] = splitted[val.Index]
				 } else if val.Type == "float" {
					 if s, err := strconv.ParseFloat(splitted[val.Index], 32); err == nil {
							m[val.Name] = math.Floor(s*100)/100
					}
				 }
				}

		}
		data, _ := json.Marshal(m)
    fmt.Printf("%s", data)
		mqtt_message = string(data)
	context.SetOutput("mqtt_message", mqtt_message)
	context.SetOutput("topic", topic)
	client.Close()
	return true, nil
}
