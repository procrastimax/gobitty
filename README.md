# gobitty
A simple terminal based calculator for bitwise operations

## Installation
In an existing Go environment just clone this repo into your go/src folder.
After this go into the folder of this repository and type `go run main.go` or `go build` into your commandline.

## Usage
There are no program parameters to pass.
Just start the executable and you are ready to go!

Currently three numerical systems are supported:
decimal, binary, hexadecimal.
All have their own syntax when entering a value:

- **decimal** 12341
- **binary**  bx1010110
- **hexadecimal** 0xDEADBEEF

Negative values and floating point numbers are not supported, so only use unsigned integer values (positive numbers).
The calculations are based on 64bit integer.

### Commands:
  - **quit** -> quits the program
  - **NUMBER as** (**dec**/**bin**/**hex**) -> converts the given number into a binary, decimal or hexadecimal one
    - f.e. `bx11 as dec` is going to return `3`
  - **EXPRESSION** -> calculates the expression with valid C-type bitwise operations
    - f.e. `(0x11|1)>>2` is going to return `4`
    - currently supported bitwise operations are: `| & >> << ~` Read more about bitwise operations [here](https://en.wikipedia.org/wiki/Bitwise_operations_in_C)
