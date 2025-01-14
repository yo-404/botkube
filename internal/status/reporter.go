package status

import (
	"context"

	"github.com/sirupsen/logrus"
)

type StatusReporter interface {
	ReportDeploymentStartup(ctx context.Context) error
	ReportDeploymentShutdown(ctx context.Context) error
	ReportDeploymentFailure(ctx context.Context, errMsg string) error
	SetResourceVersion(resourceVersion int)
}

func GetReporter(remoteCfgEnabled bool, logger logrus.FieldLogger, gql GraphQLClient, resVerClient ResVerClient, cfgVersion int) StatusReporter {
	if remoteCfgEnabled {
		return newGraphQLStatusReporter(
			logger.WithField("component", "GraphQLStatusReporter"),
			gql,
			resVerClient,
			cfgVersion,
		)
	}

	return newNoopStatusReporter()
}
