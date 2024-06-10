package mock

import (
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/mock/config"
	"github.com/gophab/gophrame/core/webservice/response"
)

func Mock() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Setting.Enabled {
			c.Next()
			return
		}

		if config.Setting.Apis != nil && len(config.Setting.Apis) > 0 {
			// 1. match api url
			var api = c.Request.Method + ":" + c.Request.RequestURI
			for path, jsonFile := range config.Setting.Apis {
				if b, _ := regexp.MatchString(path, api); b {
					// if strings.HasPrefix(api, path) {
					if pwd, err := os.Getwd(); err == nil {
						logger.Debug("Current directory: ", pwd)
					}
					content, err := os.ReadFile(jsonFile)
					if err != nil {
						logger.Error("Mock file error: ", err.Error())
						response.ErrorMessage(c, http.StatusInternalServerError, http.StatusNotFound, "Mock file error")
						return
					}
					response.ResponseWithJson(c, 200, string(content)).Abort()

					// file, err := os.Open(jsonFile)
					// if err != nil {
					// 	logger.Error("Mock file error: ", err.Error())
					// 	response.ErrorMessage(c, http.StatusInternalServerError, http.StatusNotFound, "Mock file error")
					// 	return
					// }
					// defer file.Close()

					// body := ""
					// reader := bufio.NewReader(file)
					// for {
					// 	line, err := reader.ReadString('\n')
					// 	if err == io.EOF {
					// 		break
					// 	}
					// 	body = body + line
					// }
					// response.ResponseWithJson(c, 200, body)
					return
				}
			}
		}
	}
}
