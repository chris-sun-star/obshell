package common

import (
	modelcommon "github.com/oceanbase/obshell/model/common"
)

func KVsToMap(kvs []modelcommon.KVPair) map[string]string {
	if kvs == nil {
		return nil
	}
	m := make(map[string]string)
	for _, kv := range kvs {
		m[kv.Key] = kv.Value
	}
	return m
}
