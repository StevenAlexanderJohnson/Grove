# Grove

Grove is a lightweight Go web framework designed to minimize boilerplate and simplify the process of bootstrapping web applications. It provides a set of reasonable default implementations for common web application concerns, such as routing, middleware, dependency injection, logging, and authentication, while also allowing developers to manually implement or override these components as needed.

## Vision

The vision behind Grove is to empower developers to build robust web applications quickly, without being locked into rigid conventions or heavyweight abstractions. Grove aims to strike a balance between productivity and flexibility by:

- **Reducing Boilerplate:** Providing sensible defaults for common patterns, so you can get started quickly.
- **Encouraging Manual Control:** Allowing you to override or extend any part of the framework with your own implementations.
- **Modular Design:** Organizing core features into composable modules, making it easy to use only what you need.
- **Clear Interfaces:** Defining clear interfaces for controllers, middleware, logging, and dependencies, so you can swap implementations as your application grows.
- **Secure by Default:** Including built-in support for secure authentication and request handling.

Grove is ideal for developers who want a pragmatic starting point for Go web applications, with the freedom to customize and scale as requirements evolve.

## Features

- Minimal and extensible routing with scopes and controllers
- Middleware chaining for request processing
- Pluggable logging interface with a default logger
- Simple dependency injection container
- Secure JWT authentication with JWE encryption
- Easy configuration and bootstrapping

## Getting Started

### Installation

To add Grove to your project, run:

```sh
go get github.com/StevenAlexanderJohnson/grove@v1.0.2
```

Optionally you can install the CLI to help with boilerplate generation:

```sh
go install github.com/StevenAlexanderJohnson/grove@v1.0.2
```

Then, import Grove in your Go code:

```go
import "github.com/StevenAlexanderJohnson/grove"
```

### Basic Use

Once you've added Grove to your project, you can build APIs using the `Builder Pattern` methods provided, or just use the defaults.

```go
import "github.com/StevenAlexanderJohnson/grove"

func main() {
	if err := grove.
		NewApp().
		WithRoute("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, Grove!"))
		})).Run(); err != nil {
		panic(err)
	}
}
```

## Code Generation

If you installed the Grove CLI you can initialize and generate boilerplate.
To generate a project go to the folder you want to initialize it within and run:

```sh
grove init <project_name>
```

The `project_name` will be used as the module name in `go.mod`.
It will generate a new folder inside of the directory you ran the command.

Inside is a skeleton project that you can immediately start using.
Once you're inside of that initialized project you can also use the CLI to generate boilerplate.

The CLI has two commands for generating boilerplate.
The first is `grove create resource`. This will generate the model, the repository, and service for interacting with the resource.

> It is important to note in order to use the CLI tool you have to be within a project initialized by Grove, or one that uses a similar file structure.
Because Grove CLI looks for the mod file so it knows what the module is called and can then generate import paths.

Once in the root folder of the initialized project you can run:

```sh
grove create resource <ResourceName> field1:<go_type> field2:<go_type> ...
```

It will then generate all the above mentioned files in their respective folders.


The second command that grove has for generating boilerplate is `grove create controller`.
This command will first run `grove create resource` first to generate all of the required files then it will generate a controller with CRUD methods defined.

You can then use this generated controller with ```scope.WithController()``` or ```app.WithController()```.

The command is used exactly like resource but you replace resource with controller:

```sh
grove create controller <ResourceName> field1:<go_type> field2:<go_type> ...
```

## Documentation

Grove is designed to be simple, explicit, and unopinionated. You are encouraged to read the source code and the examples provided in the `examples/` directory for real-world usage patterns.

### Key Concepts

- **Bootstrapping:** Grove helps you set up your project structure and common files, but does not hide or abstract your application logic.
- **Routing & Scopes:** Use `WithRoute` and `WithScope` to register handlers and organize your routes. Scopes can have their own middleware chains.
- **Middleware:** Middleware is just a function type. You can compose, chain, and write your own. Grove provides helpers, but you are free to implement your own logic.
- **Dependency Injection:** Grove includes a simple DI container for wiring dependencies, but you can use your own patterns if you prefer.
- **Logging:** Grove provides a pluggable logger interface. You can use the default logger or bring your own.
- **Authentication:** Secure JWE authentication is available, but you can opt out or replace it.

### Examples

See the [`examples/`](./examples/) directory for:
- Basic app setup
- Scoped routes
- Dependency injection
- Authentication

### Philosophy

Grove does not enforce strict conventions or hide details. You are always in control of your application’s structure and logic. The framework’s goal is to reduce boilerplate and help you get started, not to lock you in.

### API Reference

Full API documentation is coming soon. For now, refer to the GoDoc comments in the source code and the examples for usage guidance.

## License

This project is licensed under the MIT License.  
See the [LICENSE](LICENSE) file for details.
