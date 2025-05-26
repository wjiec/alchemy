alchemy
---
A mini-framework for quickly building HTTP & gRPC services


### Quick start
Before we begin, we need to install the CLI tool and proto plugin in your environment using the following commands:
```shell
go install github.com/wjiec/alchemy/cmd/alchemy@latest
go install github.com/wjiec/alchemy/cmd/protoc-gen-alchemy@latest
```

#### Create a project

First, create and navigate to the working directory where you want to create the project. Then, create a project using `alchemy`:
```shell
alchemy create project --repo example.com/tutorial tutorial
```


### License

alchemy is licensed under the MIT License. See [LICENSE](LICENSE) for the full license text.
