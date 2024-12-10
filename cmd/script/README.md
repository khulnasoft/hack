# Vendorless Khulnasoft

The Khulnasoft projects can be built without a vendor directory. This is a 
convenience for developers, and brings a number of benefits:

* It is easier to see the changes to the code and review them.
* It is easier to maintain the build and CI scripts, as they don't need to 
  filter out the vendor directory.
* The project doesn't look dated (the proper dependency management tools 
  are available for Go since 1.13+).
* No vendor directory means less possibility of accidentally traversing 
  into it by symlinks or scripts.

For more details and reasons for avoiding the vendor directory, see 
[khulnasoft/infra#134](https://github.com/khulnasoft/infra/issues/134).

## Status

The [khulnasoft/infra#134](https://github.com/khulnasoft/infra/issues/134) is 
ongoing effort. Currently, it is possible to use make projects vendorless, 
only if they don't use Khulnasoft nor Kubernetes code-generation tools. See the
epic issue for current status.

## Migration to a vendorless project

The following steps are required to migrate a project to be vendorless:

1. Update the `go.khulnasoft.com/hack` dependency to the latest version.
1. Update the project scripts to use the scripts inflation:
   ```patch
   -source $(dirname $0)/../vendor/go.khulnasoft.com/hack/release.sh
   +source "$(go run go.khulnasoft.com/hack/cmd/script release.sh)"
   ```
1. Update the `hack/tools.go` file to refer to the `go.khulnasoft.com/hack/cmd/script`
   tool:
   ```go
   package hack
   
   import (
   	_ "go.khulnasoft.com/hack/cmd/script"
   )
   ```
1. Remove the `vendor` directory.
1. Run `hack/update-deps.sh` to update the `go.mod` file(s).

### Examples of migrated projects

* [khulnasoft/func#1966](https://github.com/khulnasoft/func/pull/1966)
* [khulnasoft-extensions/kn-plugin-event#307](https://github.com/khulnasoft-extensions/kn-plugin-event/pull/307)
