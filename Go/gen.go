package main

//go:generate echo Hello, Go Generate!
//go:generate mkdir -p docs
//go:generate mkdir -p generated
//go:generate swag init --parseDependency --output ./docs
//go:generate echo Good Bye, Go Generate!
