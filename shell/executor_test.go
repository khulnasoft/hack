/*
Copyright 2020 The Khulnasoft Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package shell_test

import (
	"bytes"
	"testing"

	"go.khulnasoft.com/hack/shell"
)

func TestNewExecutor(t *testing.T) {
	assert := assertions{t: t}
	tests := []testcase{
		helloWorldTestCase(t),
		abortTestCase(t),
		failExampleCase(t),
		missingProjectLocationCase(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outB, errB bytes.Buffer
			tt.config.Out = &outB
			tt.config.Err = &errB
			executor := shell.NewExecutor(tt.config)
			err := tt.op(executor)
			if err != nil && !tt.wants.failed {
				t.Errorf("%s: \n got: %#v\nfailed: %#v", tt.name, err, tt.failed)
			}

			for _, wantOut := range tt.wants.outs {
				assert.Contains(outB.String(), wantOut)
			}
			for _, wantErr := range tt.wants.errs {
				assert.Contains(errB.String(), wantErr)
			}
		})
	}
}

func TestExecutorDefaults(t *testing.T) {
	assert := assertions{t: t}
	loc, err := shell.NewProjectLocation("..")
	assert.NoError(err)
	exec := shell.NewExecutor(shell.ExecutorConfig{
		ProjectLocation: loc,
	})
	err = exec.RunFunction(fn("true"))
	assert.NoError(err)
}

func helloWorldTestCase(t *testing.T) testcase {
	return testcase{
		"echo Hello, World!",
		config(t, func(cfg *shell.ExecutorConfig) {
			cfg.SkipDate = true
		}),
		func(exec shell.Executor) error {
			return exec.RunFunction(fn("echo"), "Hello, World!")
		},
		wants{
			outs: []string{
				"echo [OUT] Hello, World!",
			},
		},
	}
}

func abortTestCase(t *testing.T) testcase {
	return testcase{
		"abort function",
		config(t, func(cfg *shell.ExecutorConfig) {}),
		func(exec shell.Executor) error {
			return exec.RunFunction(fn("abort"), "reason")
		},
		wants{
			failed: true,
		},
	}
}

func failExampleCase(t *testing.T) testcase {
	return testcase{
		"fail-example.sh",
		config(t, func(cfg *shell.ExecutorConfig) {}),
		func(exec shell.Executor) error {
			return exec.RunScript(shell.Script{
				Label:      "fail-example.sh",
				ScriptPath: "shell/fail-example.sh",
			}, "expected err")
		},
		wants{
			failed: true,
			errs: []string{
				"expected err",
			},
		},
	}
}

func missingProjectLocationCase() testcase {
	return testcase{
		"missing project location",
		shell.ExecutorConfig{},
		func(exec shell.Executor) error {
			return exec.RunFunction(fn("id"))
		},
		wants{
			failed: true,
		},
	}
}

type wants struct {
	failed bool
	outs   []string
	errs   []string
}

type testcase struct {
	name   string
	config shell.ExecutorConfig
	op     func(exec shell.Executor) error
	wants
}

func config(t *testing.T, customize func(cfg *shell.ExecutorConfig)) shell.ExecutorConfig {
	assert := assertions{t: t}
	loc, err := shell.NewProjectLocation("..")
	assert.NoError(err)
	cfg := shell.ExecutorConfig{
		ProjectLocation: loc,
	}
	customize(&cfg)
	return cfg
}

func fn(name string) shell.Function {
	return shell.Function{
		Script: shell.Script{
			Label:      name,
			ScriptPath: "library.sh",
		},
		FunctionName: name,
	}
}
