package log

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lovego/xiaomei"
	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/utils"
)

var isDevMode = config.DevMode()

func Write(req *xiaomei.Request, res *xiaomei.Response, t time.Time, err interface{}) {
	fields := getFields(req, res, t)

	if err != nil {
		errStr := fmt.Sprint(err)
		errStack := string(utils.Stack(3))
		fields[`err`] = errStr
		fields[`stack`] = errStack

		if !isDevMode {
			go config.Alarm(`500错误`, formatFields(fields, false), errStr+` `+errStack)
		}
	} else if err = fields[`err`]; err != nil && !isDevMode {
		errStr, ok := err.(string)
		if !ok {
			errStr = fmt.Sprint(err)
		}
		go config.Alarm(`错误`, formatFields(fields, false), errStr)
	}

	if line := serializeFields(fields); len(line) > 0 {
		if err != nil {
			getErrorLog().Write(line)
		} else {
			getAccessLog().Write(line)
		}
	}
}

func serializeFields(fields map[string]interface{}) []byte {
	if isDevMode {
		return []byte(formatFields(fields, true))
	}
	line, err := json.Marshal(fields)
	if err != nil {
		utils.Log(`writeLog:` + err.Error())
	}
	if len(line) > 0 {
		line = append(line, '\n')
	}
	return line
}
