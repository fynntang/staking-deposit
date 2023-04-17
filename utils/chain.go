package utils

import (
	"github.com/fynntang/staking-deposit/constants"
	"strings"
)

func CheckNetworkName(networkName string) bool {
	for _, name := range constants.NetworkNames {
		if strings.Compare(string(name), networkName) == 0 {
			return true
		}
	}
	return false
}
