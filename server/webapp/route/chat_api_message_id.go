package route

import (
	"github.com/kataras/iris/context"
	"github.com/pharosnet/auid"
)

func chatApiMessageIdNew(ctx context.Context)  {
	ctx.JSON(result{Success:true, Data:map[string]interface{}{"id": auid.NewAuid()}})
	return
}
