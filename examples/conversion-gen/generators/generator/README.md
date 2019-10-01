This package exposes `ConversionGenerator`, a public struct that fulfills the `Generator` interface in an generic, extensible way to be able to write custom generators for converting similar structs across two (or more) packages.

See for example [how kubernetes wraps `ConversionGenerator` into its own `genConversion` generator struct to generate its conversion functions](https://github.com/kubernetes/gengo/blob/master/examples/conversion-gen/generators/conversion.go).
