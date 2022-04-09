package main

import (
	"dagger.io/dagger"
	"universe.dagger.io/docker"
	"universe.dagger.io/docker/cli"
  "universe.dagger.io/alpine"
  "universe.dagger.io/go"
)

dagger.#Plan & {
  client: {
    filesystem: {
      ".": read: {
        contents: dagger.#FS
        exclude: [
          "bin",
          "dev",
          "docs",
          "opni.cue",
          "testbin",
          "web/dist/_nuxt",
          "web/dist/200.html",
          "web/dist/favicon.png",
          "web/dist/loading-indicator.html",
          "plugins/example",
          "internal/cmd/testenv"
        ]
      }
    }
    network: "unix:///var/run/docker.sock": connect: dagger.#Socket
  }
  actions: {
    // Build golang image with mage installed
    builder: docker.#Build & {
      steps: [
        docker.#Pull & {
          source: "golang:1.18"
        },
        docker.#Run & {
          command: {
            name: "go"
            args: ["install", "github.com/magefile/mage@latest"]
          }
        },
      ]
    }
    // Build with mage using the builder image
    build: {
      container: go.#Container & {
        input: builder.output
        source: client.filesystem.".".read.contents
        command: {
          name: "mage"
          args: ["build"]
        }
        export: directories: "/bin": _
      }
      output: container.export.directories."/bin"
    }
    // Build the destination alpine image
    base: alpine.#Build & {
      packages: {
        "ca-certificates": _
        "curl": _
      }
    }
    // Copy the build output to the destination alpine image
    run: docker.#Build & {
      steps: [
        docker.#Copy & {
          input: base.output,
          contents: build.output,
          exclude: ["plugins/"]
          dest: "/usr/bin/opnim"
        },
        docker.#Copy & {
          // input connects to previous step's output
          contents: build.output,
          exclude: ["opnim"]
          dest: "/var/lib/opnim/"
        },
        docker.#Set & {
          config: entrypoint: ["/usr/bin/opnim"]
        },
      ]
    }
    load: cli.#Load & {
      image: run.output
      host: client.network."unix:///var/run/docker.sock".connect
      tag: "kralicky/opni-monitoring:latest"
    }
  }
}