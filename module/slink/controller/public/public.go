package public

import (
	"net/http"
	"strings"

	"github.com/gophab/gophrame/module/slink/config"
	"github.com/gophab/gophrame/module/slink/service"

	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/router"
	"github.com/gophab/gophrame/core/webservice/request"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/errors"
)

type ShortLinkPublicController struct {
	ShortLinkService *service.ShortLinkService `inject:"shortLinkService"`
}

var shortLinkPublicController = &ShortLinkPublicController{}

func Start() {
	if config.Setting.Enabled {
		contextRoot, _ := strings.CutSuffix(config.Setting.Context, "/")
		router.Root().GET(contextRoot+"/:key", shortLinkPublicController.RedirectShortLink)
	}
}

func init() {
	inject.InjectValue("shortLinkPublicController", shortLinkPublicController)
}

func (c *ShortLinkPublicController) RedirectShortLink(ctx *gin.Context) {
	key, err := request.Param(ctx, "key").MustString()
	if err != nil {
		response.FailCode(ctx, errors.INVALID_PARAMS)
		return
	}

	slink, err := c.ShortLinkService.GetByKey(key)
	if err != nil {
		response.SystemFailError(ctx, errors.MakeError("ERROR_QUERY_SHORT_LINK_ERROR"), err)
		return
	}

	if slink != nil {
		ctx.Redirect(http.StatusFound, slink.Url)
		return
	}

	response.NotFound(ctx, "Not Found or Expired")
}
