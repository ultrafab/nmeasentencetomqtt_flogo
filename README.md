# nmeasentencetomqtt_flogo
This activity convert a Nmea sentence, read the config from Redis and return a JSON-string


## Installation

```bash
flogo install github.com/ultrafab/nmeasentencetomqtt_flogo
```
Link for flogo web:
```
https://github.com/ultrafab/nmeasentencetomqtt_flogo
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "address",
      "type": "string"
    },
    {
      "name": "dbNo",
      "type": "integer"
    },
    {
      "name": "sentence",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "mqtt_message",
      "type": "string"
    },
    {
      "name": "topic",
      "type": "string"
    }
  ]
}
```
## Inputs
| Input   | Description    |
|:----------|:---------------|
| host    | the Redis address + port |
| dbNo    | the Redis database number |
| sentence    | the Nmea sentence |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| mqtt_message    | the composed mqtt messsage |
| topic    | the topic |
