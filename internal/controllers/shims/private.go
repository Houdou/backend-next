package shims

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v3"

	"github.com/penguin-statistics/backend-next/internal/constants"
	"github.com/penguin-statistics/backend-next/internal/models/cache"
	"github.com/penguin-statistics/backend-next/internal/server/svr"
	"github.com/penguin-statistics/backend-next/internal/service"
)

type PrivateController struct {
	DropMatrixService    *service.DropMatrixService
	PatternMatrixService *service.PatternMatrixService
	TrendService         *service.TrendService
	AccountService       *service.AccountService
	ItemService          *service.ItemService
	StageService         *service.StageService
}

func RegisterPrivateController(
	v2 *svr.V2,
	dropMatrixService *service.DropMatrixService,
	patternMatrixService *service.PatternMatrixService,
	trendService *service.TrendService,
	accountService *service.AccountService,
	itemService *service.ItemService,
	stageService *service.StageService,
) {
	c := &PrivateController{
		DropMatrixService:    dropMatrixService,
		PatternMatrixService: patternMatrixService,
		TrendService:         trendService,
		AccountService:       accountService,
		ItemService:          itemService,
		StageService:         stageService,
	}

	v2.Get("/_private/result/matrix/:server/:source", c.GetDropMatrix)
	v2.Get("/_private/result/pattern/:server/:source", c.GetPatternMatrix)
	v2.Get("/_private/result/trend/:server", c.GetTrends)
}

// @Summary      Get DropMatrix
// @Tags         Private
// @Produce      json
// @Param        server            path     string   true     "Server; default to CN"
// @Param        source            path     string   true     "Global or Personal; default to global"
// @Success      200               {object} shims.DropMatrixQueryResult
// @Failure      500               {object} errors.PenguinError "An unexpected error occurred"
// @Router       /PenguinStats/api/v2/_private/result/matrix/{server}/{source} [GET]
// @Deprecated
func (c *PrivateController) GetDropMatrix(ctx *fiber.Ctx) error {
	server := ctx.Params("server")
	isPersonal := ctx.Params("source") == "personal"

	accountId := null.NewInt(0, false)
	if isPersonal {
		account, err := c.AccountService.GetAccountFromRequest(ctx)
		if err != nil {
			return err
		}
		accountId.Int64 = int64(account.AccountID)
		accountId.Valid = true
	}

	shimResult, err := c.DropMatrixService.GetShimMaxAccumulableDropMatrixResults(ctx.Context(), server, true, "", "", accountId)
	if err != nil {
		return err
	}

	if !accountId.Valid {
		key := server + constants.RedisSeparator + "true"
		var lastModifiedTime time.Time
		if err := cache.LastModifiedTime.Get("[shimMaxAccumulableDropMatrixResults#server|showClosedZoned:"+key+"]", &lastModifiedTime); err != nil {
			lastModifiedTime = time.Now()
		}
		ctx.Response().Header.SetLastModified(lastModifiedTime)
	}

	return ctx.JSON(shimResult)
}

// @Summary      Get PatternMatrix
// @Tags         Private
// @Produce      json
// @Param        server            path     string   true     "Server; default to CN"
// @Param        source            path     string   true     "Global or Personal; default to global"
// @Success      200               {object} shims.PatternMatrixQueryResult
// @Failure      500               {object} errors.PenguinError "An unexpected error occurred"
// @Router       /PenguinStats/api/v2/_private/result/pattern/{server}/{source} [GET]
// @Deprecated
func (c *PrivateController) GetPatternMatrix(ctx *fiber.Ctx) error {
	server := ctx.Params("server")
	isPersonal := ctx.Params("source") == "personal"

	accountId := null.NewInt(0, false)
	if isPersonal {
		account, err := c.AccountService.GetAccountFromRequest(ctx)
		if err != nil {
			return err
		}
		accountId.Int64 = int64(account.AccountID)
		accountId.Valid = true
	}

	shimResult, err := c.PatternMatrixService.GetShimLatestPatternMatrixResults(ctx.Context(), server, accountId)
	if err != nil {
		return err
	}

	if !accountId.Valid {
		var lastModifiedTime time.Time
		if err := cache.LastModifiedTime.Get("[shimLatestPatternMatrixResults#server:"+server+"]", &lastModifiedTime); err != nil {
			lastModifiedTime = time.Now()
		}
		ctx.Response().Header.SetLastModified(lastModifiedTime)
	}

	return ctx.JSON(shimResult)
}

// @Summary      Get Trends
// @Tags         Private
// @Produce      json
// @Param        server            path     string   true     "Server; default to CN"
// @Success      200               {object} shims.TrendQueryResult
// @Failure      500               {object} errors.PenguinError "An unexpected error occurred"
// @Router       /PenguinStats/api/v2/_private/result/trend/{server} [GET]
// @Deprecated
func (c *PrivateController) GetTrends(ctx *fiber.Ctx) error {
	server := ctx.Params("server")
	shimResult, err := c.TrendService.GetShimSavedTrendResults(ctx.Context(), server)
	if err != nil {
		return err
	}

	var lastModifiedTime time.Time
	if err := cache.LastModifiedTime.Get("[shimSavedTrendResults#server:"+server+"]", &lastModifiedTime); err != nil {
		lastModifiedTime = time.Now()
	}
	ctx.Response().Header.SetLastModified(lastModifiedTime)

	return ctx.JSON(shimResult)
}
