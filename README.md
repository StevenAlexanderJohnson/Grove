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
go get github.com/StevenAlexanderJohnson/grove@v1.0.0
```

Then, import Grove in your Go code:

```go
import "github.com/StevenAlexanderJohnson/v1"
```

### Basic Use

Once you've added Grove to your project, you can build APIs using the `Builder Pattern` methods provided, or just use the defaults.

```go
import "github.com/StevenAlexanderJohnson/v1"

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

_Coming soon..._

## License

This project is licensed under the MIT License.  
See the [LICENSE](LICENSE) file for details.
