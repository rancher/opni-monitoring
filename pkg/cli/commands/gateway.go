package commands

import (
	"errors"
	"log"
	"os"

	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kralicky/opni-monitoring/pkg/auth"
	"github.com/kralicky/opni-monitoring/pkg/auth/openid"
	"github.com/kralicky/opni-monitoring/pkg/config"
	"github.com/kralicky/opni-monitoring/pkg/config/v1beta1"
	"github.com/kralicky/opni-monitoring/pkg/gateway"
	"github.com/spf13/cobra"
)

func BuildGatewayCmd() *cobra.Command {
	var configLocation string

	run := func() error {
		if configLocation == "" {
			// find config file
			path, err := config.FindConfig()
			if err != nil {
				if errors.Is(err, config.ErrConfigNotFound) {
					wd, _ := os.Getwd()
					log.Fatalf(`could not find a config file in ["%s","/etc/opni-monitoring"], and --config was not given`, wd)
				}
				log.Fatalf("an error occurred while searching for a config file: %v", err)
			}
			log.Println("using config file:", path)
			configLocation = path
		}

		objects, err := config.LoadObjectsFromFile(configLocation)
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}
		var gatewayConfig *v1beta1.GatewayConfig
		objects.Visit(
			func(config *v1beta1.GatewayConfig) {
				if gatewayConfig == nil {
					gatewayConfig = config
				}
			},
			func(ap *v1beta1.AuthProvider) {
				switch ap.Spec.Type {
				case "openid":
					mw, err := openid.New(ap.Spec)
					if err != nil {
						log.Fatalf("failed to create OpenID auth provider: %v", err)
					}
					if err := auth.InstallMiddleware(ap.GetName(), mw); err != nil {
						log.Fatalf("failed to install auth provider: %v", err)
					}
				default:
					log.Printf("unsupported auth provider type: %s", ap.Spec.Type)
				}
			},
		)

		g := gateway.NewGateway(gatewayConfig,
			gateway.WithFiberMiddleware(logger.New(), compress.New()),
			gateway.WithAuthMiddleware(gatewayConfig.Spec.AuthProvider),
			gateway.WithPrefork(false),
		)

		return g.Listen()
	}

	serveCmd := &cobra.Command{
		Use:   "gateway",
		Short: "Run the Opni Monitoring Gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			for {
				if err := run(); err != nil {
					return err
				}
			}
		},
	}

	serveCmd.Flags().StringVar(&configLocation, "config", "", "Absolute path to a config file")
	return serveCmd
}