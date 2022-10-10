package viper

import (
	"sync"

	"github.com/spf13/viper"
)

// 包一层 viper, 防止与 wundergraph pkg 的 viper 冲突
var fireBoomViper *viper.Viper
var once sync.Once

func GetFireBoomViper() *viper.Viper {
	if fireBoomViper == nil {
		once.Do(func() {
			fireBoomViper = viper.New()
		})
	}
	return fireBoomViper
}
