package main

import (
  "os"
  "context"
  
  "github.com/MontFerret/ferret/pkg/compiler"
)

func main() {
  ferret := compiler.New()
  program, err := ferret.compile(`RETURN "foobar"`)
  
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  
  out, err := program.Run(context.Background())
  
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  
  fmt.Println(string(out))
}
