# mapaccess

A small golang library to retrieve arbitrary keys from golang interface{}. Think of it as acessing JSON keys from arbitrary interfaces{}. It is heavily influenced by the golang temlpate engine.

![](https://github.com/iv-p/mapaccess/workflows/test/badge.svg)

## Installation

To install mapaccess just run the following command in your terminal
```
go get -u github.com/iv-p/mapaccess
```

## Usage

mapaccess exposes only one function, which takes a string key and a interface{}:
```
result, err := mapaccess.Get("key.one[0].two", data)
```
The key should represent a JSON type location of the data in the interface{}. It is intended to work only with pure interfaces, it won't work if you pass arbitrary structure.

## Valid keys

There is a limitation of the alphabet for keys, it includes underscores, dashes and alphanumeric characters.