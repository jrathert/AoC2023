# AoC2023
Advent of Code 2023 - in Go

## About this repository

This year I am participating [Advent of Code](https://adventofcode.com/2023/) to learn about the [Go Programming language](https://go.dev/learn/). This will be fun.

## How to use it

### Prepare for a day

There are two convenience python scripts available, to make life easier: one to download input and prepare a Go template, and 
another just to download the input

- `prepare.py` - called via `prepare.py day` it will 
    1. (try to) download the input of the specified `day`
    2. create a target directory
    3. copy the input as well as the main Go template (`main-tmpl.go`) into that directory and 
    4. start up VS Code. 
    
    E.g., `prepare.py 03` will create a director `d03` and put the respective input file `input.txt` and `main.go` into it. If the output directory already exists, the script will exit, to avoid overwriting any code. 

- `load_input.py` - called via `load_input.py day` it will 
    1. (try to) download the input of that day, 
    2. create a target directory (if it does not exist already) and 
    3. copy the input into it. 
    
    E.g., `load_input.py 03` will create a directory `d03` and put respective input file `input.txt` into that directory. If the output directory does exist, it is no error, any existing input file will be overwritten!

### Access to input files

Of course, downloading the input from the python scripts only works if the input is already available on the website (i.e. it must be at least than midnight EST/UTC-5). Also, to be able to access the input, you need to put your AoC session variable into the `.env` file - it will be read and used by the python scripts:

```
sessiontoken=abcdefgh12345678...
```

You can grab this token using your browsers development tools after logging in into Advent of Code website, see [this reddit thread](https://www.reddit.com/r/adventofcode/comments/a2vonl/how_to_download_inputs_with_a_script/).

### Happy hacking

The `prepare.py` script starts up VS Code in the main directory, but you are supposed to open the `main.go` file in the directory of the respective day. 

The `tools` directory contains a few functions that might help with everyday tasks (e.g., reading input, converting strings etc).


## License

MIT License, Copyright (c) 2023 Jonas Rathert - see file `LICENSE.txt`
