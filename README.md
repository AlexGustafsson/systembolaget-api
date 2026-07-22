# Systembolaget API

A cross-platform solution for fetching closed and open Systembolaget APIs.

The utility is an easy to use way of fetching Systembolaget's open and closed
APIs. It's written in Go and is available via several cross-platform builds.

It's usable both as a library in Go and as a standalone utility for which you
can find
[release builds here](https://github.com/AlexGustafsson/systembolaget-api/releases/).

## Quickstart

### Using as a utility

Start by grabbing the latest release for your platform from the releases or use
Docker (`ghcr.io/alexgustafsson/systembolaget-api`).

Search for beers produced in Sweden that were recently put up for sale.

```shell
systembolaget assortment --category "Öl" --origin "Sverige" --sort-by "ProductLaunchDate"
```

Get the names of Sake with a sweetness of between 5 and 12.

```shell
systembolaget assortment --taste-clock-sweetness 5-12 --category Vin --subcategory Sake | jq -cr '.[].productNameBold'
```

Get the names of non-alcoholic beverages in glass bottles.

```shell
systembolaget assortment --alcohol-percentage 0-0 --packaging-category Flaska --limit 5 | jq -cr '.[].productNameBold'
```

Download the full assortment.

```shell
systembolaget assortment --sort-by "Name" --sort ascending
```

Fetch all stores.

```shell
systembolaget stores
```

Search for a store.

```shell
systembolaget stores --search majorna
```

Get a product's status in a particular store.

```shell
# Look up the status of Guinness in the Fältöversten, Stockholm store
systembolaget stock --store-id 0102 --product-id 507849
```

An excerpt from the results is shown below. For samples, see the samples
directory.

```jsonc
{
  "productNameBold": "Melleruds",
  "productNameThin": "Utmärkta Pilsner",
  "alcoholPercentage": 4.5,
  "assortmentText": "Fast sortiment",
  "bottleText": "Flaska",
  "categoryLevel1": "Öl",
  "categoryLevel2": "Ljus lager",
  "categoryLevel3": "Pilsner - tysk stil",
  "color": "Gul färg.",
  "country": "Sverige",
  "customCategoryTitle": "Öl, Ljus lager, Pilsner - tysk stil",
  "tasteClockBitter": 6,
  "tasteClockBody": 6,
  "tasteClockCasque": 1,
  "tasteClockFruitacid": 0,
  "tasteClockSweetness": 1,
  "tasteSymbols": ["Fläsk", "Fisk", "Buffémat", "Sällskapsdryck"],
  "usage": "Serveras vid 10-12°C som sällskapsdryck, till buffé eller till rätter av fisk eller ljust kött. "
 // ...
}
```

### Using as a library

Add the necessary import.

```go
import (
 "github.com/alexgustafsson/systembolaget-api/v3/systembolaget"
)
```

Create an authenticated client.

```go
client, _ := systembolaget.DefaultClient.GetAuthenticatedClient(context.TODO())
```

Perform a search for a light lager that goes with meat.

```go
res, _ := client.Search(
 context.TODO(),
 &systembolaget.SearchOptions{
  SortBy:        systembolaget.SortPropertyScore,
  SortDirection: systembolaget.SortDirectionDescending,
 },
 systembolaget.FilterByCategory("Öl", "Ljus Lager", ""),
 systembolaget.FilterByMatch("Kött"),
)
fmt.Println(res.Products)
```

## Table of contents

[Quickstart](#quickstart)<br/>
[Use cases](#use-cases)<br/>
[Contributing](#contributing)

## Use cases

The utility can be used to automatically grab the latest available data from
Systembolagt. The data can be used to create interesting statistical charts,
archives and more. Note however that data derived from the platform should not
be used in a way that goes against
[Systembolaget's mission](https://www.omsystembolaget.se/english/systembolaget-explained/).

For archived data, please refer to <https://github.com/alexgustafsson/systembolaget-api-data>.

## Contributing

Any help with the project is more than welcome.

### Building

The CLI is easily built using Docker.

```shell
# Build for Linux using the amd64 architecture
make linux/amd64

# Run the built output
./target/linux/amd64/systembolaget
```

The docker container can also be built to run directly on your system.

```shell
docker build -t systembolaget .

# Run the container
docker run --rm -it systembolaget
```
