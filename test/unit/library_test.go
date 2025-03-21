package unit_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelperFunctions(t *testing.T) {
	t.Parallel()
	sc := newShellScript(
		fakeProwJob(),
		loadFile("source-library.bash"),
		mockGo(),
		mockKubectl(response{
			startsWith{"get pods -n test-infra --selector=app=controller"},
			simply("acme\nexample\nkhulnasoft"),
		}),
	)
	tcs := []testCase{{
		name:   `echo "$REPO_NAME_FORMATTED"`,
		stdout: lines("Khulnasoft Hack"),
	}, {
		name: `get_canonical_path test/unit/library_test.go`,
		stdout: []check{
			contains("hack/test/unit/library_test.go"),
			func(t TestingT, output string, otype outputType) bool {
				pth := strings.Trim(output, "\n")
				fi, err := os.Stat(pth)
				assert.NoError(t, err)
				assert.False(t, fi.IsDir())
				assert.True(t, path.IsAbs(pth))
				return !t.Failed()
			}},
	}, {
		name:   `capitalize "foo bar"`,
		stdout: lines("Foo Bar"),
	}, {
		name: `dump_app_logs "controller" "test-infra"`,
		stdout: lines(
			">>> Khulnasoft Hack controller logs:",
			">>> Pod: acme",
			"👻 kubectl -n test-infra logs acme --all-containers",
			">>> Pod: example",
			"👻 kubectl -n test-infra logs example --all-containers",
			">>> Pod: khulnasoft",
			"👻 kubectl -n test-infra logs khulnasoft --all-containers",
		),
	}, {
		name:   `is_protected_gcr "gcr.io/khulnasoft-releases"`,
		stdout: equal(""),
	}, {
		name:   `is_protected_gcr "gcr.io/khulnasoft-nightly"`,
		stdout: equal(""),
	}, {
		name:    `is_protected_gcr "gcr.io/khulnasoft-foobar"`,
		retcode: retcode(1),
	}, {
		name:    `is_protected_gcr "gcr.io/foobar-releases"`,
		retcode: retcode(1),
	}, {
		name:    `is_protected_gcr "gcr.io/foobar-nightly"`,
		retcode: retcode(1),
	}, {
		name:    `is_protected_gcr ""`,
		retcode: retcode(1),
	}, {
		name: `is_protected_cluster "gke_khulnasoft-tests_us-central1-f_prow"`,
	}, {
		name: `is_protected_cluster "gke_khulnasoft-tests_us-west2-a_prow"`,
	}, {
		name: `is_protected_cluster "gke_khulnasoft-tests_us-west2-a_foobar"`,
	}, {
		name:    `is_protected_cluster "gke_khulnasoft-foobar_us-west2-a_prow"`,
		retcode: retcode(1),
	}, {
		name:    `is_protected_cluster ""`,
		retcode: retcode(1),
	}, {
		name: `is_protected_project "khulnasoft-tests"`,
	}, {
		name:    `is_protected_project "khulnasoft-foobar"`,
		retcode: retcode(1),
	}, {
		name:    `is_protected_project "foobar-tests"`,
		retcode: retcode(1),
	}, {
		name:    `is_protected_project ""`,
		retcode: retcode(1),
	}, {
		name:   `calcRetcode "An example message"`,
		stdout: lines("254"),
	}, {
		name:   `calcRetcode ""`,
		stdout: lines("1"),
	}, {
		name:   `hashCode "An example message"`,
		stdout: equal("623783294"),
	}, {
		name:   `hashCode ""`,
		stdout: equal("0"),
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, tc.test(sc))
	}
}
