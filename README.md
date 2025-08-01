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
go get github.com/StevenAlexanderJohnson/grove@v0
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
