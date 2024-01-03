package netconf

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	p2pconfig "github.com/onflow/flow-go/network/p2p/config"
)

const (
	// All constant strings are used for CLI flag names and corresponding keys for config values.
	// network configuration
	networkingConnectionPruning       = "networking-connection-pruning"
	preferredUnicastsProtocols        = "preferred-unicast-protocols"
	receivedMessageCacheSize          = "received-message-cache-size"
	peerUpdateInterval                = "peerupdate-interval"
	dnsCacheTTL                       = "dns-cache-ttl"
	disallowListNotificationCacheSize = "disallow-list-notification-cache-size"
	// resource manager config
	rootResourceManagerPrefix  = "libp2p-resource-manager"
	memoryLimitRatioPrefix     = "memory-limit-ratio"
	fileDescriptorsRatioPrefix = "file-descriptors-ratio"
	limitsOverridePrefix       = "limits-override"
	systemScope                = "system"
	transientScope             = "transient"
	protocolScope              = "protocol"
	peerScope                  = "peer"
	peerProtocolScope          = "peer-protocol"
	inboundStreamLimit         = "streams-inbound"
	outboundStreamLimit        = "streams-outbound"
	inboundConnectionLimit     = "connections-inbound"
	outboundConnectionLimit    = "connections-outbound"
	fileDescriptorsLimit       = "fd"
	memoryLimitBytes           = "memory-bytes"

	alspDisabled                       = "alsp-disable-penalty"
	alspSpamRecordCacheSize            = "alsp-spam-record-cache-size"
	alspSpamRecordQueueSize            = "alsp-spam-report-queue-size"
	alspHearBeatInterval               = "alsp-heart-beat-interval"
	alspSyncEngineBatchRequestBaseProb = "alsp-sync-engine-batch-request-base-prob"
	alspSyncEngineRangeRequestBaseProb = "alsp-sync-engine-range-request-base-prob"
	alspSyncEngineSyncRequestProb      = "alsp-sync-engine-sync-request-prob"

	scoringRegistryStartupSilenceDuration = "gossipsub-scoring-registry-startup-silence-duration"
)

func AllFlagNames() []string {
	allFlags := []string{
		networkingConnectionPruning,
		preferredUnicastsProtocols,
		receivedMessageCacheSize,
		peerUpdateInterval,
		BuildFlagName(unicastKey, MessageTimeoutKey),
		BuildFlagName(unicastKey, unicastManagerKey, createStreamBackoffDelayKey),
		BuildFlagName(unicastKey, unicastManagerKey, streamZeroRetryResetThresholdKey),
		BuildFlagName(unicastKey, unicastManagerKey, maxStreamCreationRetryAttemptTimesKey),
		BuildFlagName(unicastKey, unicastManagerKey, configCacheSizeKey),
		dnsCacheTTL,
		disallowListNotificationCacheSize,
		BuildFlagName(unicastKey, rateLimiterKey, messageRateLimitKey),
		BuildFlagName(unicastKey, rateLimiterKey, BandwidthRateLimitKey),
		BuildFlagName(unicastKey, rateLimiterKey, BandwidthBurstLimitKey),
		BuildFlagName(unicastKey, rateLimiterKey, LockoutDurationKey),
		BuildFlagName(unicastKey, rateLimiterKey, DryRunKey),
		BuildFlagName(unicastKey, enableStreamProtectionKey),
		BuildFlagName(rootResourceManagerPrefix, memoryLimitRatioPrefix),
		BuildFlagName(rootResourceManagerPrefix, fileDescriptorsRatioPrefix),
		BuildFlagName(connectionManagerKey, highWatermarkKey),
		BuildFlagName(connectionManagerKey, lowWatermarkKey),
		BuildFlagName(connectionManagerKey, silencePeriodKey),
		BuildFlagName(connectionManagerKey, gracePeriodKey),
		alspDisabled,
		alspSpamRecordCacheSize,
		alspSpamRecordQueueSize,
		alspHearBeatInterval,
		alspSyncEngineBatchRequestBaseProb,
		alspSyncEngineRangeRequestBaseProb,
		alspSyncEngineSyncRequestProb,
		BuildFlagName(gossipsubKey, p2pconfig.PeerScoringEnabledKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcTracerKey, p2pconfig.LocalMeshLogIntervalKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcTracerKey, p2pconfig.ScoreTracerIntervalKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcTracerKey, p2pconfig.RPCSentTrackerCacheSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcTracerKey, p2pconfig.RPCSentTrackerQueueCacheSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcTracerKey, p2pconfig.RPCSentTrackerNumOfWorkersKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.NumberOfWorkersKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.QueueSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.ClusterPrefixedMessageConfigKey, p2pconfig.TrackerCacheSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.ClusterPrefixedMessageConfigKey, p2pconfig.TrackerCacheDecayKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.ClusterPrefixedMessageConfigKey, p2pconfig.HardThresholdKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.NotificationCacheSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.IHaveConfigKey, p2pconfig.MaxSampleSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.IHaveConfigKey, p2pconfig.MaxMessageIDSampleSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.GraftPruneMessageMaxSampleSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.IWantConfigKey, p2pconfig.MaxSampleSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.IWantConfigKey, p2pconfig.MaxMessageIDSampleSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.IWantConfigKey, p2pconfig.CacheMissThresholdKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.IWantConfigKey, p2pconfig.CacheMissCheckSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.IWantConfigKey, p2pconfig.DuplicateMsgIDThresholdKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.MessageMaxSampleSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.MessageErrorThresholdKey),
		BuildFlagName(gossipsubKey, p2pconfig.SubscriptionProviderKey, p2pconfig.UpdateIntervalKey),
		BuildFlagName(gossipsubKey, p2pconfig.SubscriptionProviderKey, p2pconfig.CacheSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.AppSpecificScoreRegistryKey, p2pconfig.ScoreUpdateWorkerNumKey),
		BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.AppSpecificScoreRegistryKey, p2pconfig.ScoreUpdateRequestQueueSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.AppSpecificScoreRegistryKey, p2pconfig.ScoreTTLKey),
		BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.SpamRecordCacheKey, p2pconfig.CacheSizeKey),
		BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.SpamRecordCacheKey, p2pconfig.PenaltyDecaySlowdownThresholdKey),
		BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.SpamRecordCacheKey, p2pconfig.DecayRateReductionFactorKey),
		BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.SpamRecordCacheKey, p2pconfig.PenaltyDecayEvaluationPeriodKey),
		BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.DecayIntervalKey),
		BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.ScoringRegistryStartupSilenceDurationKey),
	}

	for _, scope := range []string{systemScope, transientScope, protocolScope, peerScope, peerProtocolScope} {
		for _, resource := range []string{inboundStreamLimit,
			outboundStreamLimit,
			inboundConnectionLimit,
			outboundConnectionLimit,
			fileDescriptorsLimit,
			memoryLimitBytes} {
			allFlags = append(allFlags, fmt.Sprintf("%s-%s-%s-%s", rootResourceManagerPrefix, limitsOverridePrefix, scope, resource))
		}
	}

	return allFlags
}

// InitializeNetworkFlags initializes all CLI flags for the Flow network configuration on the provided pflag set.
// Args:
//
//	*pflag.FlagSet: the pflag set of the Flow node.
//	*Config: the default network config used to set default values on the flags
func InitializeNetworkFlags(flags *pflag.FlagSet, config *Config) {
	flags.Bool(networkingConnectionPruning, config.NetworkConnectionPruning, "enabling connection trimming")
	flags.Duration(dnsCacheTTL, config.DNSCacheTTL, "time-to-live for dns cache")
	flags.StringSlice(
		preferredUnicastsProtocols, config.PreferredUnicastProtocols, "preferred unicast protocols in ascending order of preference")
	flags.Uint32(receivedMessageCacheSize, config.NetworkReceivedMessageCacheSize, "incoming message cache size at networking layer")
	flags.Uint32(
		disallowListNotificationCacheSize,
		config.DisallowListNotificationCacheSize,
		"cache size for notification events from disallow list")
	flags.Duration(peerUpdateInterval, config.PeerUpdateInterval, "how often to refresh the peer connections for the node")
	flags.Duration(BuildFlagName(unicastKey, MessageTimeoutKey), config.Unicast.MessageTimeout, "how long a unicast transmission can take to complete")
	flags.Duration(BuildFlagName(unicastKey, unicastManagerKey, createStreamBackoffDelayKey), config.Unicast.UnicastManager.CreateStreamBackoffDelay,
		"initial backoff delay between failing to establish a connection with another node and retrying, "+
			"this delay increases exponentially with the number of subsequent failures to establish a connection.")
	flags.Uint64(BuildFlagName(unicastKey, unicastManagerKey, streamZeroRetryResetThresholdKey), config.Unicast.UnicastManager.StreamZeroRetryResetThreshold,
		"reset stream creation retry budget from zero to the maximum after consecutive successful streams reach this threshold.")
	flags.Uint64(BuildFlagName(unicastKey, unicastManagerKey, maxStreamCreationRetryAttemptTimesKey),
		config.Unicast.UnicastManager.MaxStreamCreationRetryAttemptTimes,
		"max attempts to create a unicast stream.")
	flags.Uint32(BuildFlagName(unicastKey, unicastManagerKey, configCacheSizeKey), config.Unicast.UnicastManager.ConfigCacheSize,
		"cache size of the dial config cache, recommended to be big enough to accommodate the entire nodes in the network.")

	// unicast stream handler rate limits
	flags.Int(BuildFlagName(unicastKey, rateLimiterKey, messageRateLimitKey), config.Unicast.RateLimiter.MessageRateLimit, "maximum number of unicast messages that a peer can send per second")
	flags.Int(BuildFlagName(unicastKey, rateLimiterKey, BandwidthRateLimitKey), config.Unicast.RateLimiter.BandwidthRateLimit,
		"bandwidth size in bytes a peer is allowed to send via unicast streams per second")
	flags.Int(BuildFlagName(unicastKey, rateLimiterKey, BandwidthBurstLimitKey), config.Unicast.RateLimiter.BandwidthBurstLimit, "bandwidth size in bytes a peer is allowed to send at one time")
	flags.Duration(BuildFlagName(unicastKey, rateLimiterKey, LockoutDurationKey), config.Unicast.RateLimiter.LockoutDuration,
		"the number of seconds a peer will be forced to wait before being allowed to successful reconnect to the node after being rate limited")
	flags.Bool(BuildFlagName(unicastKey, rateLimiterKey, DryRunKey), config.Unicast.RateLimiter.DryRun, "disable peer disconnects and connections gating when rate limiting peers")
	flags.Bool(BuildFlagName(unicastKey, enableStreamProtectionKey),
		config.Unicast.EnableStreamProtection,
		"enable stream protection for unicast streams, when enabled, all connections that are being established or have been already established for unicast streams are protected")

	LoadLibP2PResourceManagerFlags(flags, config)

	flags.Int(BuildFlagName(connectionManagerKey, lowWatermarkKey), config.ConnectionManager.LowWatermark, "low watermarking for libp2p connection manager")
	flags.Int(BuildFlagName(connectionManagerKey, highWatermarkKey), config.ConnectionManager.HighWatermark, "high watermarking for libp2p connection manager")
	flags.Duration(BuildFlagName(connectionManagerKey, gracePeriodKey), config.ConnectionManager.GracePeriod, "grace period for libp2p connection manager")
	flags.Duration(BuildFlagName(connectionManagerKey, silencePeriodKey), config.ConnectionManager.SilencePeriod, "silence period for libp2p connection manager")
	flags.Bool(BuildFlagName(gossipsubKey, p2pconfig.PeerScoringEnabledKey), config.GossipSub.PeerScoringEnabled, "enabling peer scoring on pubsub network")
	flags.Duration(BuildFlagName(gossipsubKey, p2pconfig.RpcTracerKey, p2pconfig.LocalMeshLogIntervalKey),
		config.GossipSub.RpcTracer.LocalMeshLogInterval,
		"logging interval for local mesh in gossipsub tracer")
	flags.Duration(BuildFlagName(gossipsubKey, p2pconfig.RpcTracerKey, p2pconfig.ScoreTracerIntervalKey), config.GossipSub.RpcTracer.ScoreTracerInterval,
		"logging interval for peer score tracer in gossipsub, set to 0 to disable")
	flags.Uint32(BuildFlagName(gossipsubKey, p2pconfig.RpcTracerKey, p2pconfig.RPCSentTrackerCacheSizeKey), config.GossipSub.RpcTracer.RPCSentTrackerCacheSize,
		"cache size of the rpc sent tracker used by the gossipsub mesh tracer.")
	flags.Uint32(BuildFlagName(gossipsubKey, p2pconfig.RpcTracerKey, p2pconfig.RPCSentTrackerQueueCacheSizeKey), config.GossipSub.RpcTracer.RPCSentTrackerQueueCacheSize,
		"cache size of the rpc sent tracker worker queue.")
	flags.Int(BuildFlagName(gossipsubKey, p2pconfig.RpcTracerKey, p2pconfig.RPCSentTrackerNumOfWorkersKey), config.GossipSub.RpcTracer.RpcSentTrackerNumOfWorkers,
		"number of workers for the rpc sent tracker worker pool.")
	// gossipsub RPC control message validation limits used for validation configuration and rate limiting
	flags.Int(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.NumberOfWorkersKey),
		config.GossipSub.RpcInspector.Validation.NumberOfWorkers,
		"number of gossipsub RPC control message validation inspector component workers")
	flags.Uint32(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.QueueSizeKey),
		config.GossipSub.RpcInspector.Validation.QueueSize,
		"queue size for gossipsub RPC validation inspector events worker pool queue.")
	flags.Uint32(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.ClusterPrefixedMessageConfigKey, p2pconfig.TrackerCacheSizeKey),
		config.GossipSub.RpcInspector.Validation.ClusterPrefixedMessage.ControlMsgsReceivedCacheSize,
		"cache size for gossipsub RPC validation inspector cluster prefix received tracker.")
	flags.Float64(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.ClusterPrefixedMessageConfigKey, p2pconfig.TrackerCacheDecayKey),
		config.GossipSub.RpcInspector.Validation.ClusterPrefixedMessage.ControlMsgsReceivedCacheDecay,
		"the decay value used to decay cluster prefix received topics received cached counters.")
	flags.Float64(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.ClusterPrefixedMessageConfigKey, p2pconfig.HardThresholdKey),
		config.GossipSub.RpcInspector.Validation.ClusterPrefixedMessage.HardThreshold,
		"the maximum number of cluster-prefixed control messages allowed to be processed when the active cluster id is unset or a mismatch is detected, exceeding this threshold will result in node penalization by gossipsub.")
	// networking event notifications
	flags.Uint32(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.NotificationCacheSizeKey), config.GossipSub.RpcInspector.NotificationCacheSize,
		"cache size for notification events from gossipsub rpc inspector")
	// application layer spam prevention (alsp) protocol
	flags.Bool(alspDisabled, config.AlspConfig.DisablePenalty, "disable the penalty mechanism of the alsp protocol. default value (recommended) is false")
	flags.Uint32(alspSpamRecordCacheSize, config.AlspConfig.SpamRecordCacheSize, "size of spam record cache, recommended to be 10x the number of authorized nodes")
	flags.Uint32(alspSpamRecordQueueSize, config.AlspConfig.SpamReportQueueSize, "size of spam report queue, recommended to be 100x the number of authorized nodes")
	flags.Duration(alspHearBeatInterval,
		config.AlspConfig.HearBeatInterval,
		"interval between two consecutive heartbeat events at alsp, recommended to leave it as default unless you know what you are doing.")
	flags.Float32(alspSyncEngineBatchRequestBaseProb,
		config.AlspConfig.SyncEngine.BatchRequestBaseProb,
		"base probability of creating a misbehavior report for a batch request message")
	flags.Float32(alspSyncEngineRangeRequestBaseProb,
		config.AlspConfig.SyncEngine.RangeRequestBaseProb,
		"base probability of creating a misbehavior report for a range request message")
	flags.Float32(alspSyncEngineSyncRequestProb, config.AlspConfig.SyncEngine.SyncRequestProb, "probability of creating a misbehavior report for a sync request message")

	flags.Float64(BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.SpamRecordCacheKey, p2pconfig.PenaltyDecaySlowdownThresholdKey),
		config.GossipSub.ScoringParameters.SpamRecordCache.PenaltyDecaySlowdownThreshold,
		fmt.Sprintf("the penalty level at which the decay rate is reduced by --%s",
			BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.SpamRecordCacheKey, p2pconfig.DecayRateReductionFactorKey)))
	flags.Float64(BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.SpamRecordCacheKey, p2pconfig.DecayRateReductionFactorKey),
		config.GossipSub.ScoringParameters.SpamRecordCache.DecayRateReductionFactor,
		fmt.Sprintf("defines the value by which the decay rate is decreased every time the penalty is below the --%s",
			BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.SpamRecordCacheKey, p2pconfig.PenaltyDecaySlowdownThresholdKey)))
	flags.Duration(BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.SpamRecordCacheKey, p2pconfig.PenaltyDecayEvaluationPeriodKey),
		config.GossipSub.ScoringParameters.SpamRecordCache.PenaltyDecayEvaluationPeriod,
		"defines the period at which the decay for a spam record is okay to be adjusted.")
	flags.Int(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.IHaveConfigKey, p2pconfig.MaxSampleSizeKey),
		config.GossipSub.RpcInspector.Validation.IHave.MaxSampleSize,
		"max number of ihaves to sample when performing validation")
	flags.Int(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.IHaveConfigKey, p2pconfig.MaxMessageIDSampleSizeKey),
		config.GossipSub.RpcInspector.Validation.IHave.MaxMessageIDSampleSize,
		"max number of message ids to sample when performing validation per ihave")
	flags.Int(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.GraftPruneMessageMaxSampleSizeKey),
		config.GossipSub.RpcInspector.Validation.GraftPruneMessageMaxSampleSize,
		"max number of control messages to sample when performing validation on GRAFT and PRUNE message types")
	flags.Uint(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.IWantConfigKey, p2pconfig.MaxSampleSizeKey),
		config.GossipSub.RpcInspector.Validation.IWant.MaxSampleSize,
		"max number of iwants to sample when performing validation")
	flags.Int(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.IWantConfigKey, p2pconfig.MaxMessageIDSampleSizeKey),
		config.GossipSub.RpcInspector.Validation.IWant.MaxMessageIDSampleSize,
		"max number of message ids to sample when performing validation per iwant")
	flags.Float64(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.IWantConfigKey, p2pconfig.CacheMissThresholdKey),
		config.GossipSub.RpcInspector.Validation.IWant.CacheMissThreshold,
		"max number of iwants to sample when performing validation")
	flags.Int(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.IWantConfigKey, p2pconfig.CacheMissCheckSizeKey),
		config.GossipSub.RpcInspector.Validation.IWant.CacheMissCheckSize,
		"the iWants size at which message id cache misses will be checked")
	flags.Float64(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.IWantConfigKey, p2pconfig.DuplicateMsgIDThresholdKey),
		config.GossipSub.RpcInspector.Validation.IWant.DuplicateMsgIDThreshold,
		"max allowed duplicate message IDs in a single iWant control message")
	flags.Int(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.MessageMaxSampleSizeKey),
		config.GossipSub.RpcInspector.Validation.MessageMaxSampleSize,
		"the max sample size used for RPC message validation. If the total number of RPC messages exceeds this value a sample will be taken but messages will not be truncated")
	flags.Int(BuildFlagName(gossipsubKey, p2pconfig.RpcInspectorKey, p2pconfig.ValidationConfigKey, p2pconfig.MessageErrorThresholdKey),
		config.GossipSub.RpcInspector.Validation.MessageErrorThreshold,
		"the threshold at which an error will be returned if the number of invalid RPC messages exceeds this value")
	flags.Duration(BuildFlagName(gossipsubKey, p2pconfig.SubscriptionProviderKey, p2pconfig.UpdateIntervalKey),
		config.GossipSub.SubscriptionProvider.UpdateInterval,
		"interval for updating the list of subscribed topics for all peers in the gossipsub, recommended value is a few minutes")
	flags.Uint32(BuildFlagName(gossipsubKey, p2pconfig.SubscriptionProviderKey, p2pconfig.CacheSizeKey),
		config.GossipSub.SubscriptionProvider.CacheSize,
		"size of the cache that keeps the list of topics each peer has subscribed to, recommended size is 10x the number of authorized nodes")

	flags.Duration(BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.ScoringRegistryStartupSilenceDurationKey),
		config.GossipSub.ScoringParameters.ScoringRegistryStartupSilenceDuration,
		"the duration of time, after the node startup, during which the scoring registry remains inactive before penalizing nodes.")
	flags.Int(BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.AppSpecificScoreRegistryKey, p2pconfig.ScoreUpdateWorkerNumKey),
		config.GossipSub.ScoringParameters.AppSpecificScore.ScoreUpdateWorkerNum,
		"number of workers for the app specific score update worker pool")
	flags.Uint32(BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.AppSpecificScoreRegistryKey, p2pconfig.ScoreUpdateRequestQueueSizeKey),
		config.GossipSub.ScoringParameters.AppSpecificScore.ScoreUpdateRequestQueueSize,
		"size of the app specific score update worker pool queue")
	flags.Duration(BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.AppSpecificScoreRegistryKey, p2pconfig.ScoreTTLKey),
		config.GossipSub.ScoringParameters.AppSpecificScore.ScoreTTL,
		"time to live for app specific scores; when expired a new request will be sent to the score update worker pool; till then the expired score will be used")
	flags.Uint32(BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.SpamRecordCacheKey, p2pconfig.CacheSizeKey),
		config.GossipSub.ScoringParameters.SpamRecordCache.CacheSize,
		"size of the spam record cache, recommended size is 10x the number of authorized nodes")
	flags.Duration(BuildFlagName(gossipsubKey, p2pconfig.ScoreParamsKey, p2pconfig.DecayIntervalKey),
		config.GossipSub.ScoringParameters.DecayInterval,
		"interval at which the counters associated with a peer behavior in GossipSub system are decayed, recommended value is one minute")
}

// LoadLibP2PResourceManagerFlags loads all CLI flags for the libp2p resource manager configuration on the provided pflag set.
// Args:
// *pflag.FlagSet: the pflag set of the Flow node.
// *Config: the default network config used to set default values on the flags
func LoadLibP2PResourceManagerFlags(flags *pflag.FlagSet, config *Config) {
	flags.Float64(fmt.Sprintf("%s-%s", rootResourceManagerPrefix, fileDescriptorsRatioPrefix),
		config.ResourceManager.FileDescriptorsRatio,
		"ratio of available file descriptors to be used by libp2p (in (0,1])")
	flags.Float64(fmt.Sprintf("%s-%s", rootResourceManagerPrefix, memoryLimitRatioPrefix),
		config.ResourceManager.MemoryLimitRatio,
		"ratio of available memory to be used by libp2p (in (0,1])")
	loadLibP2PResourceManagerFlagsForScope(systemScope, flags, &config.ResourceManager.Override.System)
	loadLibP2PResourceManagerFlagsForScope(transientScope, flags, &config.ResourceManager.Override.Transient)
	loadLibP2PResourceManagerFlagsForScope(protocolScope, flags, &config.ResourceManager.Override.Protocol)
	loadLibP2PResourceManagerFlagsForScope(peerScope, flags, &config.ResourceManager.Override.Peer)
	loadLibP2PResourceManagerFlagsForScope(peerProtocolScope, flags, &config.ResourceManager.Override.PeerProtocol)
}

// loadLibP2PResourceManagerFlagsForScope loads all CLI flags for the libp2p resource manager configuration on the provided pflag set for the specific scope.
// Args:
// *p2pconf.ResourceScope: the resource scope to load flags for.
// *pflag.FlagSet: the pflag set of the Flow node.
// *Config: the default network config used to set default values on the flags.
func loadLibP2PResourceManagerFlagsForScope(scope p2pconfig.ResourceScope, flags *pflag.FlagSet, override *p2pconfig.ResourceManagerOverrideLimit) {
	flags.Int(fmt.Sprintf("%s-%s-%s-%s", rootResourceManagerPrefix, limitsOverridePrefix, scope, inboundStreamLimit),
		override.StreamsInbound,
		fmt.Sprintf("the limit on the number of inbound streams at %s scope, 0 means use the default value", scope))
	flags.Int(fmt.Sprintf("%s-%s-%s-%s", rootResourceManagerPrefix, limitsOverridePrefix, scope, outboundStreamLimit),
		override.StreamsOutbound,
		fmt.Sprintf("the limit on the number of outbound streams at %s scope, 0 means use the default value", scope))
	flags.Int(fmt.Sprintf("%s-%s-%s-%s", rootResourceManagerPrefix, limitsOverridePrefix, scope, inboundConnectionLimit),
		override.ConnectionsInbound,
		fmt.Sprintf("the limit on the number of inbound connections at %s scope, 0 means use the default value", scope))
	flags.Int(fmt.Sprintf("%s-%s-%s-%s", rootResourceManagerPrefix, limitsOverridePrefix, scope, outboundConnectionLimit),
		override.ConnectionsOutbound,
		fmt.Sprintf("the limit on the number of outbound connections at %s scope, 0 means use the default value", scope))
	flags.Int(fmt.Sprintf("%s-%s-%s-%s", rootResourceManagerPrefix, limitsOverridePrefix, scope, fileDescriptorsLimit),
		override.FD,
		fmt.Sprintf("the limit on the number of file descriptors at %s scope, 0 means use the default value", scope))
	flags.Int(fmt.Sprintf("%s-%s-%s-%s", rootResourceManagerPrefix, limitsOverridePrefix, scope, memoryLimitBytes),
		override.Memory,
		fmt.Sprintf("the limit on the amount of memory (bytes) at %s scope, 0 means use the default value", scope))
}

// SetAliases this func sets an aliases for each CLI flag defined for network config overrides to it's corresponding
// full key in the viper config store. This is required because in our p2pconfig.yml file all configuration values for the
// Flow network are stored one level down on the network-config property. When the default config is bootstrapped viper will
// store these values with the "network-p2pconfig." prefix on the config key, because we do not want to use CLI flags like --network-p2pconfig.networking-connection-pruning
// to override default values we instead use cleans flags like --networking-connection-pruning and create an alias from networking-connection-pruning -> network-p2pconfig.networking-connection-pruning
// to ensure overrides happen as expected.
// Args:
// *viper.Viper: instance of the viper store to register network config aliases on.
// Returns:
// error: if a flag does not have a corresponding key in the viper store; all returned errors are fatal.
func SetAliases(conf *viper.Viper) error {
	m := make(map[string]string)
	// create map of key -> full pathkey
	// ie: "networking-connection-pruning" -> "network-p2pconfig.networking-connection-pruning"
	for _, key := range conf.AllKeys() {
		s := strings.Split(key, ".")
		// Each networking config has the format of network-p2pconfig.key1.key2.key3... in the config file
		// which is translated to key1-key2-key3... in the CLI flags
		// Hence, we map the CLI flag name to the full key in the config store
		// TODO: all networking flags should also be prefixed with "network-config". Hence, this
		// mapping should be from network-p2pconfig.key1.key2.key3... to network-config-key1-key2-key3...
		m[strings.Join(s[1:], "-")] = key
	}
	// each flag name should correspond to exactly one key in our config store after it is loaded with the default config
	for _, flagName := range AllFlagNames() {
		fullKey, ok := m[flagName]
		if !ok {
			return fmt.Errorf("invalid network configuration missing configuration key flag name %s check config file and cli flags", flagName)
		}
		conf.RegisterAlias(fullKey, flagName)
	}
	return nil
}

func BuildFlagName(keys ...string) string {
	return strings.Join(keys, "-")
}
