// +build unit

package documents

import (
	"os"
	"testing"

	"github.com/centrifuge/go-centrifuge/config/configstore"

	"github.com/centrifuge/go-centrifuge/bootstrap"
	"github.com/centrifuge/go-centrifuge/bootstrap/bootstrappers/testlogging"
	"github.com/centrifuge/go-centrifuge/config"
	"github.com/centrifuge/go-centrifuge/queue"
	"github.com/centrifuge/go-centrifuge/storage"
	"github.com/centrifuge/go-centrifuge/transactions"
)

var ctx = map[string]interface{}{}
var ConfigService configstore.Service

func TestMain(m *testing.M) {
	ibootstappers := []bootstrap.TestBootstrapper{
		&testlogging.TestLoggingBootstrapper{},
		&config.Bootstrapper{},
		&storage.Bootstrapper{},
		&configstore.Bootstrapper{},
		transactions.Bootstrapper{},
		&queue.Bootstrapper{},
		&Bootstrapper{},
	}
	bootstrap.RunTestBootstrappers(ibootstappers, ctx)
	ConfigService = ctx[configstore.BootstrappedConfigStorage].(configstore.Service)
	result := m.Run()
	bootstrap.RunTestTeardown(ibootstappers)
	os.Exit(result)
}