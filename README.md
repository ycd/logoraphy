

<p align="center"><img src="/examples/logoraphy.jpeg" align="center" width="736" height="413"></div>

## A typographic logo generator.
------

## Installation

#### From GitHub

```
git clone https://github.com/ycd/logoraphy.git
cd logoraphy
go build .
./logoraphy -h
```

## Usage

### Usage with command line interface

```go
Usage of ./logoraphy:
  -bg string
        Background color.
  -name string
        Company name, a string.
  -type string
        Output format, JPEG or PNG. (default "jpeg")
```

### Serving as API endpoint

```
git switch fiber-serving
go run .
```

### Usage from api

```
Path params:
    companyName string

Query params:
    type: string
        - Default: jpeg.
        - jpeg or png
    bg: string
        - Default: Generates random color.
        - indigo, red, primary, gray..
    
   
```
## Examples

Netflix              |  Facebook
:-------------------------:|:-------------------------:
![](/examples/netflix.jpeg)  |  ![](/examples/facebook.jpeg)
