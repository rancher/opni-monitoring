package agent

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/prometheus/model/rulefmt"
	"github.com/rancher/opni-monitoring/pkg/rules"
	"github.com/rancher/opni-monitoring/pkg/sdk/api"
	"github.com/rancher/opni-monitoring/pkg/util"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func (a *Agent) configureRuleFinder() (rules.RuleFinder, error) {
	if a.Rules != nil {
		if pr := a.Rules.Discovery.PrometheusRules; pr != nil {
			client, err := util.NewK8sClient(util.ClientOptions{
				Kubeconfig: pr.Kubeconfig,
				Scheme:     api.NewScheme(),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create k8s client: %w", err)
			}
			finder := rules.NewPrometheusRuleFinder(client,
				rules.WithLogger(a.logger),
				rules.WithNamespaces(pr.SearchNamespaces...),
			)
			return finder, nil
		}
	}
	return nil, fmt.Errorf("missing configuration")
}

func (a *Agent) streamRuleGroupUpdates(ctx context.Context) (<-chan [][]byte, error) {
	finder, err := a.configureRuleFinder()
	if err != nil {
		return nil, fmt.Errorf("failed to configure rule discovery: %w", err)
	}
	searchInterval := time.Minute * 15
	if interval := a.Rules.Discovery.Interval; interval != "" {
		duration, err := time.ParseDuration(interval)
		if err != nil {
			return nil, fmt.Errorf("failed to parse discovery interval: %w", err)
		}
		searchInterval = duration
	}
	notifier := rules.NewPeriodicUpdateNotifier(ctx, finder, searchInterval)

	notifierC := notifier.NotifyC(ctx)
	a.logger.Debug("starting rule group update notifier")
	groupYamlDocs := make(chan [][]byte, cap(notifierC))
	go func() {
		defer close(groupYamlDocs)
		for {
			ruleGroups, ok := <-notifierC
			if !ok {
				a.logger.Debug("rule discovery channel closed")
				return
			}
			a.logger.Debug("received updated rule groups from discovery")
			go func() {
				groupYamlDocs <- a.marshalRuleGroups(ruleGroups)
			}()
		}
	}()
	return groupYamlDocs, nil
}

func (a *Agent) marshalRuleGroups(ruleGroups []rulefmt.RuleGroup) [][]byte {
	yamlDocs := make([][]byte, 0, len(ruleGroups))
	for _, ruleGroup := range ruleGroups {
		doc, err := yaml.Marshal(ruleGroup)
		if err != nil {
			a.logger.With(
				zap.Error(err),
				zap.String("group", ruleGroup.Name),
			).Error("failed to marshal rule group")
			continue
		}
		yamlDocs = append(yamlDocs, doc)
	}
	return yamlDocs
}

func (a *Agent) streamRulesToGateway(ctx context.Context) error {
	lg := a.logger
	updateC, err := a.streamRuleGroupUpdates(ctx)
	if err != nil {
		a.logger.With(
			zap.Error(err),
		).Error("failed to configure rule discovery")
		return err
	}
	pending := make(chan [][]byte, 1)
	go func() {
		defer close(pending)
		for {
			var docs [][]byte
			select {
			case docs = <-pending:
			case <-ctx.Done():
				lg.Error(ctx.Err())
				return
			}
		RETRY:
			lg.Debug("sending alert rules to gateway")
			for {
				for _, doc := range docs {
					reqCtx, ca := context.WithTimeout(ctx, time.Second*2)
					defer ca()
					code, _, err := a.gatewayClient.Post(reqCtx, "/api/agent/sync_rules").
						Set("Content-Type", "application/yaml").
						Body(doc).
						Do()
					if err != nil || code != http.StatusAccepted {
						// retry, unless another update is received from the channel
						lg.With(
							zap.Error(err),
							zap.Int("code", code),
						).Error("failed to send alert rules to gateway (retry in 5 seconds)")
						select {
						case docs = <-pending:
							lg.Debug("updated rules were received during backoff, retrying immediately")
							goto RETRY
						case <-time.After(5 * time.Second):
							goto RETRY
						case <-ctx.Done():
							lg.Error(ctx.Err())
							return
						}
					}
				}
				lg.Infof("successfully sent %d alert rules to gateway", len(docs))
				break
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		case yamlDocs, ok := <-updateC:
			if !ok {
				lg.Debug("rule discovery stream closed")
				return nil
			}
			pending <- yamlDocs
		}
	}
}
