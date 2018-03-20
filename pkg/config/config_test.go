package config

import (
	"testing"
)

func TestUnmarshal(t *testing.T) {
	_, err := Load([]byte(`a B`))
	if err == nil {
		t.Errorf("Expect an error")
	}
	data := `
name: test_backup
logfile: /var/log/test_backup.log`
	c, _ := Load([]byte(data))
	if c.Name != "test_backup" {
		t.Errorf("Uncorrect unmarched values")
	}

}
