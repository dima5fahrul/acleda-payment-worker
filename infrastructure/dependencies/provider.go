package dependencies

import (
	"payment-airpay/application/services"
	"payment-airpay/infrastructure/database"
	"payment-airpay/infrastructure/database/connectors"
	"payment-airpay/infrastructure/database/repositories"
	"payment-airpay/infrastructure/gateway/acleda"
	"payment-airpay/infrastructure/publishers"
	"payment-airpay/infrastructure/service"
	"sync"

	"github.com/google/wire"
)

// singleton
var gatewayOnce sync.Once
var transactionServiceOnce sync.Once
var publisherOnce sync.Once
var yugabyteClientOnce sync.Once
var masterDataRepoOnce sync.Once
var paymentRepoOnce sync.Once
var acledaRepoOnce sync.Once

// singleton instance
var acledaGatewayInstance *acleda.AcledaGateway
var transactionServiceInstance *service.PaymentAcleda
var publisherInstance *publishers.PublisherLog
var yugabyteClientInstance *connectors.YugabyteConnector
var masterDataRepoInstance *repositories.MasterDataRepositoryYugabyteDB
var paymentRepoInstance *repositories.PaymentRepositoryYugabyteDB
var acledaRepoInstance *repositories.AcledaRepositoryYugabyteDB

var ProviderSet wire.ProviderSet = wire.NewSet(
	ProvideAcledaGateway,
	ProvideTransactionService,
	ProvideYugabyteClient,
	ProvideMasterDataRepository,
	ProvidePaymentRepository,
	ProvideAcledaRepository,
	ProvidePublisher,
	wire.Bind(new(services.PaymentGateway), new(*acleda.AcledaGateway)),
	wire.Bind(new(services.TransactionService), new(*service.PaymentAcleda)),
	wire.Bind(new(services.Publisher), new(*publishers.PublisherLog)),
)

func ProvideAcledaGateway() *acleda.AcledaGateway {
	gatewayOnce.Do(func() {
		acledaGatewayInstance = acleda.NewAcledaGateway()
	})
	return acledaGatewayInstance
}

func ProvideTransactionService() *service.PaymentAcleda {
	transactionServiceOnce.Do(func() {
		masterRepo := ProvideMasterDataRepository()
		paymentRepo := ProvidePaymentRepository()
		acledaRepo := ProvideAcledaRepository()
		db := ProvideYugabyteClient()
		transactionServiceInstance = service.NewPaymentAcleda(masterRepo, paymentRepo, acledaRepo, db)
	})
	return transactionServiceInstance
}

func ProvideYugabyteClient() *connectors.YugabyteConnector {
	yugabyteClientOnce.Do(func() {
		yugabyteClientInstance = connectors.NewYugabyteConnector(database.YugabyteDBClient)
	})
	return yugabyteClientInstance
}

func ProvideMasterDataRepository() *repositories.MasterDataRepositoryYugabyteDB {
	masterDataRepoOnce.Do(func() {
		masterDataRepoInstance = repositories.NewMasterDataRepositoryYugabyteDB()
	})
	return masterDataRepoInstance
}

func ProvidePaymentRepository() *repositories.PaymentRepositoryYugabyteDB {
	paymentRepoOnce.Do(func() {
		paymentRepoInstance = repositories.NewPaymentRepositoryYugabyteDB()
	})
	return paymentRepoInstance
}

func ProvideAcledaRepository() *repositories.AcledaRepositoryYugabyteDB {
	acledaRepoOnce.Do(func() {
		acledaRepoInstance = repositories.NewAcledaRepositoryYugabyteDB()
	})
	return acledaRepoInstance
}

func ProvidePaymentAcledaService() *service.PaymentAcleda {
	return service.NewPaymentAcleda(
		ProvideMasterDataRepository(),
		ProvidePaymentRepository(),
		ProvideAcledaRepository(),
		ProvideYugabyteClient(),
	)
}

func ProvidePublisher() *publishers.PublisherLog {
	publisherOnce.Do(func() {
		publisherInstance = publishers.NewPublisherLog()
	})
	return publisherInstance
}
