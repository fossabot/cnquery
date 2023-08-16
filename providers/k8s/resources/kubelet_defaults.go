/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// https://github.com/kubernetes/kubernetes/blob/release-1.25/pkg/kubelet/apis/config/v1beta1/defaults.go
package resources

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"

	// TODO: Cut references to k8s.io/kubernetes, eventually there should be none from this package
	logsapi "k8s.io/component-base/logs/api/v1"
	utilpointer "k8s.io/utils/pointer"
)

const (
	// TODO: Move these constants to k8s.io/kubelet/config/v1beta1 instead?
	DefaultIPTablesMasqueradeBit = 14
	DefaultIPTablesDropBit       = 15
	DefaultVolumePluginDir       = "/usr/libexec/kubernetes/kubelet-plugins/volume/exec/"

	// See https://github.com/kubernetes/enhancements/tree/master/keps/sig-node/2570-memory-qos
	DefaultMemoryThrottlingFactor = 0.8
)

// https://github.com/kubernetes/kubernetes/blob/release-1.25/pkg/kubelet/types/constants.go
const (
	// ResolvConfDefault is the system default DNS resolver configuration.
	ResolvConfDefault = "/etc/resolv.conf"
)

// https://github.com/kubernetes/kubernetes/blob/release-1.25/pkg/cluster/ports/ports.go
const (
	// KubeletPort is the default port for the kubelet server on each host machine.
	// May be overridden by a flag at startup.
	KubeletPort = 10250
	// KubeletReadOnlyPort exposes basic read-only services from the kubelet.
	// May be overridden by a flag at startup.
	// This is necessary for heapster to collect monitoring stats from the kubelet
	// until heapster can transition to using the SSL endpoint.
	// TODO(roberthbailey): Remove this once we have a better solution for heapster.
	KubeletReadOnlyPort = 10255
	// KubeletHealthzPort exposes a healthz endpoint from the kubelet.
	// May be overridden by a flag at startup.
	KubeletHealthzPort = 10248
)

// https://github.com/kubernetes/kubernetes/blob/release-1.25/pkg/kubelet/qos/policy.go
const (
	// KubeletOOMScoreAdj is the OOM score adjustment for Kubelet
	KubeletOOMScoreAdj int = -999
)

// https://github.com/kubernetes/kubernetes/blob/release-1.25/pkg/kubelet/apis/config/v1beta1/defaults_linux.go
// DefaultEvictionHard includes default options for hard eviction.
var DefaultEvictionHard = map[string]string{
	"memory.available":  "100Mi",
	"nodefs.available":  "10%",
	"nodefs.inodesFree": "5%",
	"imagefs.available": "15%",
}

var (
	zeroDuration = metav1.Duration{}
	// TODO: Move these constants to k8s.io/kubelet/config/v1beta1 instead?
	// Refer to [Node Allocatable](https://git.k8s.io/design-proposals-archive/node/node-allocatable.md) doc for more information.
	DefaultNodeAllocatableEnforcement = []string{"pods"}
)

func SetDefaults_KubeletConfiguration(obj *kubeletconfigv1beta1.KubeletConfiguration) {
	if obj.EnableServer == nil {
		obj.EnableServer = utilpointer.BoolPtr(true)
	}
	if obj.SyncFrequency == zeroDuration {
		obj.SyncFrequency = metav1.Duration{Duration: 1 * time.Minute}
	}
	if obj.FileCheckFrequency == zeroDuration {
		obj.FileCheckFrequency = metav1.Duration{Duration: 20 * time.Second}
	}
	if obj.HTTPCheckFrequency == zeroDuration {
		obj.HTTPCheckFrequency = metav1.Duration{Duration: 20 * time.Second}
	}
	if obj.Address == "" {
		obj.Address = "0.0.0.0"
	}
	if obj.Port == 0 {
		obj.Port = KubeletPort
	}
	if obj.Authentication.Anonymous.Enabled == nil {
		obj.Authentication.Anonymous.Enabled = utilpointer.BoolPtr(false)
	}
	if obj.Authentication.Webhook.Enabled == nil {
		obj.Authentication.Webhook.Enabled = utilpointer.BoolPtr(true)
	}
	if obj.Authentication.Webhook.CacheTTL == zeroDuration {
		obj.Authentication.Webhook.CacheTTL = metav1.Duration{Duration: 2 * time.Minute}
	}
	if obj.Authorization.Mode == "" {
		obj.Authorization.Mode = kubeletconfigv1beta1.KubeletAuthorizationModeWebhook
	}
	if obj.Authorization.Webhook.CacheAuthorizedTTL == zeroDuration {
		obj.Authorization.Webhook.CacheAuthorizedTTL = metav1.Duration{Duration: 5 * time.Minute}
	}
	if obj.Authorization.Webhook.CacheUnauthorizedTTL == zeroDuration {
		obj.Authorization.Webhook.CacheUnauthorizedTTL = metav1.Duration{Duration: 30 * time.Second}
	}
	if obj.RegistryPullQPS == nil {
		obj.RegistryPullQPS = utilpointer.Int32Ptr(5)
	}
	if obj.RegistryBurst == 0 {
		obj.RegistryBurst = 10
	}
	if obj.EventRecordQPS == nil {
		obj.EventRecordQPS = utilpointer.Int32Ptr(5)
	}
	if obj.EventBurst == 0 {
		obj.EventBurst = 10
	}
	if obj.EnableDebuggingHandlers == nil {
		obj.EnableDebuggingHandlers = utilpointer.BoolPtr(true)
	}
	if obj.HealthzPort == nil {
		obj.HealthzPort = utilpointer.Int32Ptr(10248)
	}
	if obj.HealthzBindAddress == "" {
		obj.HealthzBindAddress = "127.0.0.1"
	}
	if obj.OOMScoreAdj == nil {
		obj.OOMScoreAdj = utilpointer.Int32Ptr(int32(KubeletOOMScoreAdj))
	}
	if obj.StreamingConnectionIdleTimeout == zeroDuration {
		obj.StreamingConnectionIdleTimeout = metav1.Duration{Duration: 4 * time.Hour}
	}
	if obj.NodeStatusReportFrequency == zeroDuration {
		// For backward compatibility, NodeStatusReportFrequency's default value is
		// set to NodeStatusUpdateFrequency if NodeStatusUpdateFrequency is set
		// explicitly.
		if obj.NodeStatusUpdateFrequency == zeroDuration {
			obj.NodeStatusReportFrequency = metav1.Duration{Duration: 5 * time.Minute}
		} else {
			obj.NodeStatusReportFrequency = obj.NodeStatusUpdateFrequency
		}
	}
	if obj.NodeStatusUpdateFrequency == zeroDuration {
		obj.NodeStatusUpdateFrequency = metav1.Duration{Duration: 10 * time.Second}
	}
	if obj.NodeLeaseDurationSeconds == 0 {
		obj.NodeLeaseDurationSeconds = 40
	}
	if obj.ImageMinimumGCAge == zeroDuration {
		obj.ImageMinimumGCAge = metav1.Duration{Duration: 2 * time.Minute}
	}
	if obj.ImageGCHighThresholdPercent == nil {
		// default is below docker's default dm.min_free_space of 90%
		obj.ImageGCHighThresholdPercent = utilpointer.Int32Ptr(85)
	}
	if obj.ImageGCLowThresholdPercent == nil {
		obj.ImageGCLowThresholdPercent = utilpointer.Int32Ptr(80)
	}
	if obj.VolumeStatsAggPeriod == zeroDuration {
		obj.VolumeStatsAggPeriod = metav1.Duration{Duration: time.Minute}
	}
	if obj.CgroupsPerQOS == nil {
		obj.CgroupsPerQOS = utilpointer.BoolPtr(true)
	}
	if obj.CgroupDriver == "" {
		obj.CgroupDriver = "cgroupfs"
	}
	if obj.CPUManagerPolicy == "" {
		obj.CPUManagerPolicy = "none"
	}
	if obj.CPUManagerReconcilePeriod == zeroDuration {
		// Keep the same as default NodeStatusUpdateFrequency
		obj.CPUManagerReconcilePeriod = metav1.Duration{Duration: 10 * time.Second}
	}
	if obj.MemoryManagerPolicy == "" {
		obj.MemoryManagerPolicy = kubeletconfigv1beta1.NoneMemoryManagerPolicy
	}
	if obj.TopologyManagerPolicy == "" {
		obj.TopologyManagerPolicy = kubeletconfigv1beta1.NoneTopologyManagerPolicy
	}
	if obj.TopologyManagerScope == "" {
		obj.TopologyManagerScope = kubeletconfigv1beta1.ContainerTopologyManagerScope
	}
	if obj.RuntimeRequestTimeout == zeroDuration {
		obj.RuntimeRequestTimeout = metav1.Duration{Duration: 2 * time.Minute}
	}
	if obj.HairpinMode == "" {
		obj.HairpinMode = kubeletconfigv1beta1.PromiscuousBridge
	}
	if obj.MaxPods == 0 {
		obj.MaxPods = 110
	}
	// default nil or negative value to -1 (implies node allocatable pid limit)
	if obj.PodPidsLimit == nil || *obj.PodPidsLimit < int64(0) {
		obj.PodPidsLimit = utilpointer.Int64(-1)
	}

	if obj.ResolverConfig == nil {
		obj.ResolverConfig = utilpointer.String(ResolvConfDefault)
	}
	if obj.CPUCFSQuota == nil {
		obj.CPUCFSQuota = utilpointer.BoolPtr(true)
	}
	if obj.CPUCFSQuotaPeriod == nil {
		obj.CPUCFSQuotaPeriod = &metav1.Duration{Duration: 100 * time.Millisecond}
	}
	if obj.NodeStatusMaxImages == nil {
		obj.NodeStatusMaxImages = utilpointer.Int32Ptr(50)
	}
	if obj.MaxOpenFiles == 0 {
		obj.MaxOpenFiles = 1000000
	}
	if obj.ContentType == "" {
		obj.ContentType = "application/vnd.kubernetes.protobuf"
	}
	if obj.KubeAPIQPS == nil {
		obj.KubeAPIQPS = utilpointer.Int32Ptr(5)
	}
	if obj.KubeAPIBurst == 0 {
		obj.KubeAPIBurst = 10
	}
	if obj.SerializeImagePulls == nil {
		obj.SerializeImagePulls = utilpointer.BoolPtr(true)
	}
	if obj.EvictionHard == nil {
		obj.EvictionHard = DefaultEvictionHard
	}
	if obj.EvictionPressureTransitionPeriod == zeroDuration {
		obj.EvictionPressureTransitionPeriod = metav1.Duration{Duration: 5 * time.Minute}
	}
	if obj.EnableControllerAttachDetach == nil {
		obj.EnableControllerAttachDetach = utilpointer.BoolPtr(true)
	}
	if obj.MakeIPTablesUtilChains == nil {
		obj.MakeIPTablesUtilChains = utilpointer.BoolPtr(true)
	}
	if obj.IPTablesMasqueradeBit == nil {
		obj.IPTablesMasqueradeBit = utilpointer.Int32Ptr(DefaultIPTablesMasqueradeBit)
	}
	if obj.IPTablesDropBit == nil {
		obj.IPTablesDropBit = utilpointer.Int32Ptr(DefaultIPTablesDropBit)
	}
	if obj.FailSwapOn == nil {
		obj.FailSwapOn = utilpointer.BoolPtr(true)
	}
	if obj.ContainerLogMaxSize == "" {
		obj.ContainerLogMaxSize = "10Mi"
	}
	if obj.ContainerLogMaxFiles == nil {
		obj.ContainerLogMaxFiles = utilpointer.Int32Ptr(5)
	}
	if obj.ConfigMapAndSecretChangeDetectionStrategy == "" {
		obj.ConfigMapAndSecretChangeDetectionStrategy = kubeletconfigv1beta1.WatchChangeDetectionStrategy
	}
	if obj.EnforceNodeAllocatable == nil {
		obj.EnforceNodeAllocatable = DefaultNodeAllocatableEnforcement
	}
	if obj.VolumePluginDir == "" {
		obj.VolumePluginDir = DefaultVolumePluginDir
	}
	// Use the Default LoggingConfiguration option
	logsapi.SetRecommendedLoggingConfiguration(&obj.Logging)
	if obj.EnableSystemLogHandler == nil {
		obj.EnableSystemLogHandler = utilpointer.BoolPtr(true)
	}
	if obj.EnableProfilingHandler == nil {
		obj.EnableProfilingHandler = utilpointer.BoolPtr(true)
	}
	if obj.EnableDebugFlagsHandler == nil {
		obj.EnableDebugFlagsHandler = utilpointer.BoolPtr(true)
	}
	if obj.SeccompDefault == nil {
		obj.SeccompDefault = utilpointer.BoolPtr(false)
	}
	if obj.MemoryThrottlingFactor == nil {
		obj.MemoryThrottlingFactor = utilpointer.Float64Ptr(DefaultMemoryThrottlingFactor)
	}
	if obj.RegisterNode == nil {
		obj.RegisterNode = utilpointer.BoolPtr(true)
	}
	if obj.LocalStorageCapacityIsolation == nil {
		obj.LocalStorageCapacityIsolation = utilpointer.BoolPtr(true)
	}
}
