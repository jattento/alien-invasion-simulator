[![Last release](https://img.shields.io/github/v/release/jattento/alien-invasion-simulator?style=plastic)](https://github.com/jattento/alien-invasion-simulator/releases)
[![Build Status](https://github.com/jattento/alien-invasion-simulator/actions/workflows/go.yml/badge.svg)](https://github.com/jattento/alien-invasion-simulator/actions/workflows/go.yml)

# alien-invasion-simulator
Simulate an alien invasion. You know... just in case.

![](https://static.displate.com/280x392/displate/2020-12-21/c7eda619a5ddfadde906a368aaa7f212_4590ac66ec1dcfc80ae53403da627140.jpg)

## Usage

![](https://drive.google.com/uc?export=view&id=1xfk999PCTae-QVUVhgKnAxwUmneevSBA)

To download the program, just clone the repo or get the compiled code from the [releases page](https://github.com/jattento/alien-invasion-simulator/releases)

First time simulating the earth's fall? 🔥 The easiest way to run this for the first time is to just
run `alien-sim` in your console.

Nice, now that you are an advanced user 👨‍🔬, 
and probably want to simulate the outcome for the cities of your country:
you can load the custom cities layout using the `--city-config=path` flag
where the path is should be pointing to a valid text file in your file system.

Oh, you just want to see some random world burn? 😈 The system can create the city's layout file for you
use the flags `--matrix` and `--cities` to indicate how it should look like.

```
Usage:
    alien-sim [flags]

Flags:
    -a, --aliens int            Amount of aliens to spawn (default 15)
    -c, --cities int            Amount of cities deployed in the matrix (default 20)
        --city-config string    Path where to find the city config file.
    -d, --days int              Days until simulation ends. (default 10000)
    -m, --matrix int            Matrix size where the value is N when N*N=total matrix size. (default 5)
```

Also keep in mind the controls used inside the simulation:

- `Control + Q`: Close
- `Control + A`: Time speed down
- `Control + S`: Time speed up

## Config file format
Example:
```
Foo north=Bar west=Baz south=Qu-ux
Bar south=Foo west=Bee
Qu-ux north=Foo
Baz east=Foo
Bee east=Bar
```
Rules:
- Each city that appears in the file must have its own unique record
- Layout must be consistent: if Baz has Foo at the east, then Foo must have Baz at the west
- The city and each of the pairs are separated by a single space, and the
  directions are separated from their respective cities with an equals (=) sign.

## Scaffolding
This repo was designed using [package oriented design](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html).

## Code Style
This repo follows the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

## Considerations

- No third party graph package was utilized since this functionality is core to the exercise
and using an existing solutions would be a waste of an opportunity to show off data structure knowledge

## Assumptions

- All the aliens travel at the same time and cannot fight each other at that moment.