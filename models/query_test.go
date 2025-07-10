package models

import (
	"testing"
)

func TestExpandVariable(t *testing.T) {
	CurrentEnv = &Environment{
		Variables: map[string]string{
			"host": "free.fr",
			"id":   "12",
		},
	}
	res, err := ExpandMapVariable(map[string]string{
		"test":  "test",
		"test2": "http://%{host}%/api/v1/%{id}%",
	})
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}

	if res["test"] != "test" {
		t.Errorf("test should be test not %s\n", res["test"])
		return
	}
	if res["test2"] != "http://free.fr/api/v1/12" {
		t.Errorf("test2 should be http://free.fr/api/v1/12 not %s\n", res["test2"])
		return
	}

}
