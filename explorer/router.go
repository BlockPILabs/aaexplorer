package explorer

import (
	"github.com/BlockPILabs/aa-scan/explorer/controller"
	"github.com/gofiber/fiber/v2"
	fiber_recover "github.com/gofiber/fiber/v2/middleware/recover"
)

func Resister(router fiber.Router) {
	router.Use(fiber_recover.New())
	v1 := router.Group("/v1")
	// networks

	v1.Get("/networks", controller.GetNetworks).Name(controller.NameGetNetworks)

	networksV1 := v1.Group("/network/:network<regex(^[a-z0-9]{1,}$)}>?")
	networksV1.Use(controller.NetworkMiddleware())
	networksV1.Get("/", controller.GetNetwork)

	// search
	networksV1.Get("/search", controller.SearchAll).Name(controller.NameSearchAll)

	// Bundlers
	networksV1.Get("/bundlers", controller.GetBundlers).Name(controller.NameGetBundlers)
	networksV1.Get("/bundler/:bundler<regex(0x[a-z0-9]{40}$)}>", controller.GetBundler).Name(controller.NameGetBundler)

	// bundles
	networksV1.Get("/bundles", controller.GetBundles).Name(controller.NameGetBundles)
	networksV1.Get("/bundle/:bundle<regex(0x[a-z0-9]{40}$)}>", controller.GetBundle).Name(controller.NameGetBundle)
	// factory
	networksV1.Get("/factories", controller.GetFactories).Name(controller.NameGetFactories)

	// userops
	networksV1.Get("/userops", controller.GetUserOps).Name(controller.NameGetUserOps)
	networksV1.Get("/useropsAnalysis", controller.GetUserOpsAnalysis).Name(controller.NameGetUserOpsAnalysis)
	networksV1.Get("/useropsAnalysisList", controller.GetUserOpsAnalysisList).Name(controller.NameGetUserOpsAnalysisList)

	// paymasters
	networksV1.Get("/paymasters", controller.GetPaymasters).Name(controller.NameGetPaymasters)
	// blocks
	networksV1.Get("/blocks", controller.GetBlocks).Name(controller.NameGetBlocks)
	networksV1.Get("/aaBlocks", controller.GetAABlocksPage).Name(controller.NameGetAABlocksPage)
	networksV1.Get("/block/:block<regex((^0x[a-z0-9]{64}$|^\\d+$))}>", controller.GetBlock).Name(controller.NameGetBlock)

	//transaction
	networksV1.Get("/transaction", controller.GetTransaction).Name(controller.NameGetTransaction)
	networksV1.Get("/transactionList", controller.GetTransactionList).Name(controller.NameGetTransactionList)

	//aaAccount
	networksV1.Get("/aaAccountInfo", controller.GetAaAccountInfo).Name(controller.NameGetAaAcountInfo)

	//user
	//networksV1.Get("/user")
	//networksV1.Get("/user")
	networksV1.Get("/userBalance", controller.GetUserBalance).Name(controller.NameGetUserBalance)

	//home page
	networksV1.Get("/dailyStatistic", controller.GetDailyStatistic).Name(controller.NameGetDailyStatistic)
	networksV1.Get("/aaTxnDominance", controller.GetAATxnDominance).Name(controller.NameGetAATxnDominance)

	//user op type analyze
	networksV1.Get("/userOpType", controller.GetUserOpType).Name(controller.NameGetUserOpType)
	networksV1.Get("/aaContractInteract", controller.GetAAContractInteract).Name(controller.NameGetAAContractInteract)

	//top list
	networksV1.Get("/topBundler", controller.GetTopBundler).Name(controller.NameGetTopBundler)
	networksV1.Get("/topPaymaster", controller.GetTopPaymaster).Name(controller.NameGetTopPaymaster)
	networksV1.Get("/topFactory", controller.GetTopFactory).Name(controller.NameGetTopFactory)
}

func Error() {

}
