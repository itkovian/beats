package gpfs

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/elastic/beats/libbeat/common"
)

// QuotaInfo contains the information of a single entry produced by mmrepquota
type QuotaInfo struct {
	filesystem string
	fileset    string
	kind       string
	entity     string
	blockUsage int64
	blockSoft  int64
	blockHard  int64
	blockDoubt int64
	blockGrace string
	filesUsage int64
	filesSoft  int64
	filesHard  int64
	filesDoubt int64
	filesGrace string
}

// MmRepQuota is a wrapper around the mmrepquota command
func MmRepQuota() ([]QuotaInfo, error) {

	cmd := exec.Command("mmrepquota") // TODO: pass arguments
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		var nope []QuotaInfo
		return nope, errors.New("mmrepquota failed")
	}

	var quotas []QuotaInfo
	quotas, err = parseMmRepQuota(out.String())
	if err != nil {
		var nope []QuotaInfo
		return nope, errors.New("mmrepquota info could not be parsed")
	}
	return quotas, nil
}

// GetQuotaEvent turns the quota information into a MapStr
func GetQuotaEvent(quota *QuotaInfo) common.MapStr {
	return common.MapStr{
		"filesystem":    quota.filesystem,
		"fileset":       quota.fileset,
		"kind":          quota.kind,
		"entity":        quota.entity,
		"block_usage":   quota.blockUsage,
		"block_soft":    quota.blockSoft,
		"block_hard":    quota.blockHard,
		"block_doubt":   quota.blockDoubt,
		"block_expired": quota.blockGrace,
		"files_usage":   quota.filesUsage,
		"files_soft":    quota.filesSoft,
		"files_hard":    quota.filesHard,
		"files_doubt":   quota.filesDoubt,
		"files_expired": quota.filesGrace,
	}
}

// ParseMmRepQuota converts the lines into the desired information
func parseMmRepQuota(output string) (qs []QuotaInfo, err error) {

	lines := strings.Split(output, "\n")
	fieldMap := parseMmRepQuotaHeader(lines[0])

	for _, line := range lines[1:] {
		fields := strings.Split(line, ":")
		qi := QuotaInfo{
			filesystem: fields[fieldMap["filesystemName"]],
			fileset:    fields[fieldMap["filesetname"]],
			kind:       fields[fieldMap["quotaType"]],
			entity:     fields[fieldMap["name"]],
			blockUsage: strconv.ParseInt(fields[fieldMap["blockUsage"]], 64),
			blockSoft:  strconv.ParseInt(fields[fieldMap["blockQuota"]], 64),
			blockHard:  strconv.ParseInt(fields[fieldMap["blockLimit"]], 64),
			blockDoubt: strconv.ParseInt(fields[fieldMap["blockInDoubt"]], 64),
			blockGrace: fields[fieldMap["blockGrace"]],
			filesUsage: strconv.ParseInt(fields[fieldMap["filesUsage"]], 64),
			filesSoft:  strconv.ParseInt(fields[fieldMap["filesQuota"]], 64),
			filesHard:  strconv.ParseInt(fields[fieldMap["filesLimit"]], 64),
			filesDoubt: strconv.ParseInt(fields[fieldMap["filesInDoubt"]], 64),
			filesGrace: fields[fieldMap["filesGrace"]],
		}
		qs.append(qi)
	}
	return
}

// parseMmRepQuotaHeader builds a map of the field names and the corresponding index
func parseMmRepQuotaHeader(header string) (m map[string]int) {

	for i, s := range strings.Split(header, ":") {
		if s != "" {
			continue
		}
		m[s] = i
	}

	return
}
