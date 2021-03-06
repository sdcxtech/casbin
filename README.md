# casbin

Another [casbin](https://casbin.org/) implementation in golang.

## Difference with the official casbin implementation

* Use google [Common Expression Language](https://github.com/google/cel-go) as the matcher expression language.
* Assertion field in policy and request only can be `string` type. So there is no support for `ABAC` model.
* Only implement the core feature checking permissions. Not include policies and roles management
  which should be implemented in a different [Bounded Context](https://martinfowler.com/bliki/BoundedContext.html).

## License

Released under the [MIT License](LICENSE).
