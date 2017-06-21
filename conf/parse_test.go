package conf

import (
	"testing"
)

func TestGetConfWithoutParse(t *testing.T) {
	_, err := GetConf()
	if err != ErrConfigNotInit {
		t.Error("Error should return if get conf without parse")
	}
}

func TestParseConfWithNotExistFile(t *testing.T) {
	_, err := ParseConf("/not/exist/file")
	if err == nil {
		t.Error("Error should return if read file from not exist file")
	}
}
