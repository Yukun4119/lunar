package common

import (
	"github.com/stretchr/testify/assert"
	"lunar_uml/util"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	filePath := "../../config/config.yml"
	config := util.LoadConfig(filePath)
	assert.NotNilf(t, config, "not nil")
}
