/*
Copyright 2021 The Khulnasoft Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"

	"go.khulnasoft.com/hack/schema/commands"
	"go.khulnasoft.com/hack/schema/example"
	"go.khulnasoft.com/hack/schema/registry"
)

// This is a demo of what the CLI looks like, copy and implement your own.
func main() {
	registry.Register(&example.LoremIpsum{})

	if err := commands.New("go.khulnasoft.com/hack/schema").Execute(); err != nil {
		log.Fatal("Error during command execution: ", err)
	}
}
