package quota

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/metricbeat/mb"
	"github.com/elastic/beats/metricbeat/module/gpfs"
)

var debugf = logp.MakeDebug("gpfs-quota")

// init registers the MetricSet with the central registry.
// The New method will be called after the setup of the module and before starting to fetch data
func init() {
	if err := mb.Registry.AddMetricSet("gpfs", "quota", New); err != nil {
		logp.Err("Aaaaargh, no cigar")
		panic(err)
	}
	logp.Warn("quota init ran")
}

// MetricSet type defines all fields of the MetricSet
// As a minimum it must inherit the mb.BaseMetricSet fields, but can be extended with
// additional entries. These variables can be used to persist data or configuration between
// multiple fetch calls.
type MetricSet struct {
	mb.BaseMetricSet
	Filesystem string
}

// New create a new instance of the MetricSet
// Part of new is also setting up the configuration by processing additional
// configuration entries if needed.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	config := struct {
		Filesystem string `config:"filesystems"`
	}{
		Filesystem: "",
	}

	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, err
	}
	logp.Warn("Creating a new instance of the quota Metricset")

	return &MetricSet{
		BaseMetricSet: base,
		Filesystem:    config.Filesystem,
	}, nil
}

// Fetch methods implements the data gathering and data conversion to the right format
// It returns the event which is then forward to the output. In case of an error, a
// descriptive error must be returned.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {

	gpfsQuota, err := gpfs.MmRepQuota(m.Filesystem) // TODO: get this for each filesystem
	logp.Warn("retrieved quota information from mmrepquota")
	if err != nil {
		panic("Could not get quota information")
	}

	quota := make([]common.MapStr, 0, len(gpfsQuota))
	for _, q := range gpfsQuota {
		quota = append(quota, gpfs.GetQuotaEvent(&q))
	}

	return quota, nil
}
