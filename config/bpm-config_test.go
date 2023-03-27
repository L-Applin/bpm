package config

import (
	"reflect"
	"testing"
)

func TestMergeEmptyConfigs(t *testing.T) {
	assertConfigs(t, Config{}, Config{}, Config{})
	assertConfigs(t, Config{}, Config{"key1": "value1"}, Config{"key1": "value1"})
	assertConfigs(t, Config{"key1": "value1"}, Config{}, Config{"key1": "value1"})
}

func TestMergeConfigs(t *testing.T) {
	assertConfigs(t,
		Config{
			"key-1": "value-1",
		}, Config{
			"key-2": "value-2",
		}, Config{
			"key-1": "value-1",
			"key-2": "value-2",
		})

	assertConfigs(t,
		Config{
			"key-1": "value-1",
			"key-2": "value-2",
		},
		Config{
			"key-2": "new-value-2",
			"key-3": "value-3",
		},
		Config{
			"key-1": "value-1",
			"key-2": "new-value-2",
			"key-3": "value-3",
		})

	assertConfigs(t,
		Config{
			"key-1": "value-1",
			"key-2": "value-2",
			"key-3": "value-3",
		},
		Config{
			"key-1": "new-value-1",
			"key-2": "new-value-2",
			"key-3": "new-value-3",
		},
		Config{
			"key-1": "new-value-1",
			"key-2": "new-value-2",
			"key-3": "new-value-3",
		})
}

func TestMergeNestedConfigs(t *testing.T) {
	assertConfigs(t,
		Config{
			"key-1": "value-1",
			"key-2": Config{
				"nested-key-1": "value-2",
				"nested-key-2": "value-3",
			},
		},
		Config{
			"key-2": Config{
				"nested-key-2": "new-value-3",
				"nested-key-3": "value-4",
			},
		},
		Config{
			"key-1": "value-1",
			"key-2": Config{
				"nested-key-1": "value-2",
				"nested-key-2": "new-value-3",
				"nested-key-3": "value-4",
			},
		})

	assertConfigs(t,
		Config{
			"key-0": "value-0",
			"key-1": "value-1",
			"key-2": Config{
				"nested-key-0": "nested-value-0",
				"nested-key-1": "nested-value-1",
				"nested-key-2": Config{
					"nested-nested-key-0": "nested-nested-value-0",
					"nested-nested-key-1": "nested-nested-value-1",
				},
				"nested-key-3": Config{
					"nested-nested-key-0.1": "nested-nested-value-0.1",
					"nested-nested-key-1.1": "nested-nested-value-1.1",
				},
			},
		},
		Config{
			"key-1": "new-value-1",
			"key-2": Config{
				"nested-key-1": "new-nested-value-1",
				"nested-key-2": Config{
					"nested-nested-key-1": "new-nested-nested-value-1",
					"nested-nested-key-2": "new-nested-nested-value-2",
					"nested-nested-key-3": Config{
						"nested-nested-nested-key-4": "value-6",
					},
				},
				"nested-key-4": Config{
					"nested-nested-key-1.1": "new-nested-nested-values-1.1",
				},
			},
		},
		Config{
			"key-0": "value-0",
			"key-1": "new-value-1",
			"key-2": Config{
				"nested-key-0": "nested-value-0",
				"nested-key-1": "new-nested-value-1",
				"nested-key-2": Config{
					"nested-nested-key-0": "nested-nested-value-0",
					"nested-nested-key-1": "new-nested-nested-value-1",
					"nested-nested-key-2": "new-nested-nested-value-2",
					"nested-nested-key-3": Config{
						"nested-nested-nested-key-4": "value-6",
					},
				},
				"nested-key-3": Config{
					"nested-nested-key-0.1": "nested-nested-value-0.1",
					"nested-nested-key-1.1": "nested-nested-value-1.1",
				},
				"nested-key-4": Config{
					"nested-nested-key-1.1": "new-nested-nested-values-1.1",
				},
			},
		})

}

func assertConfigs(t *testing.T, first, second, expected Config) {
	t.Helper()
	assertConfigsEquals(t, MergeConfigs(first, second), expected)
}

func assertConfigsEquals(t *testing.T, actual, expected Config) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected configs to be equals. \nExpected: %#v\nActual:%#v\n", expected, actual)
	}
}
