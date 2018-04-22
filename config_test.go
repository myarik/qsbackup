package qsbackup

import (
	"testing"
)

func TestUnmarshal(t *testing.T) {
	_, err := ConfigLoad([]byte(`a B`))
	if err == nil {
		t.Errorf("Expect an error")
	}
//	data := `
//name: test_backup
//logfile: /var/log/test_backup.log`
//	c, _ := ConfigLoad([]byte(data))
//	if c.Name != "test_backup" {
//		t.Errorf("Uncorrect unmarched values")
//	}

}
