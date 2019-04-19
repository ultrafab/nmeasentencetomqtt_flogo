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
	Func string `json:"func"`
	Indexes []struct {
				Index int `json:"index"`
				Name string `json:"name"`
				Type string `json:"type"`
		} `json:"indexes"`
}

func MyFunction(a, b int) int {
  return a + b
}

func ConvertCoordinates(stringa string, dir string) float32{
	s_seconds := string(stringa[len(stringa)-3:len(stringa)-1])
	s_minutes := string(stringa[len(stringa)-6:len(stringa)-4])
	s_degrees := string(stringa[0:len(stringa)-6])
	seconds,errs := strconv.Atoi(s_seconds)
	if errs != nil {
		log.Errorf("Error! Could not parse : %s", s_seconds )
	}
	minutes,errm := strconv.Atoi(s_minutes)
	if errm != nil {
		log.Errorf("Error! Could not parse : %s", s_minutes )
	}
	degrees,errd := strconv.Atoi(s_degrees)
	if errd != nil {
		log.Errorf("Error! Could not parse : %s", s_degrees )
	}
	cseconds := float32(seconds)*float32(60)/float32(100)
	v := float32(degrees) + float32(minutes)/float32(60) + float32(cseconds)/float32(3600)
	if dir == "S" || dir == "W"{
		v = v * float32(-1)
	}
	return v
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
				if v.Func == "convertCoordinates" {
					lat_string := fmt.Sprintf("%.3f", m["lat"])
					lat_dir := fmt.Sprintf("%s",m["lat_dir"])
					lat_decimal := ConvertCoordinates(lat_string,lat_dir)
					m["lat_decimal"] = lat_decimal
					lon_string := fmt.Sprintf("%.3f", m["lon"])
					lon_dir := fmt.Sprintf("%s",m["lon_dir"])
					lon_decimal := ConvertCoordinates(lon_string,lon_dir)
					m["lon_decimal"] = lon_decimal
				}
		}
		data, _ := json.Marshal(m)
		mqtt_message = string(data)
	context.SetOutput("mqtt_message", mqtt_message)
	context.SetOutput("topic", topic)
	client.Close()
	return true, nil
}
