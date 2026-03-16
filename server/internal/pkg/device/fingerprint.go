package device

import (
	"crypto/sha256"
	"encoding/hex"
)

// GenerateFingerprint 生成设备指纹
// 基于传入的设备信息生成唯一指纹
func GenerateFingerprint(cpuID, motherboardSN, macAddress string) string {
	data := cpuID + "|" + motherboardSN + "|" + macAddress
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}