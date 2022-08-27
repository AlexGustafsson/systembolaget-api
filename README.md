systembolaget-api  - a cross-platform solution for fetching and converting [Systembolaget's open APIs](https://www.systembolaget.se/api/) to JSON
======

The utility is an easy to use way of fetching Systembolaget's open APIs and converting them to JSON and XML. It's written in Go and is available via several cross-platform builds.

It's usable both as a library in Go and as a standalone utility for which you can find release builds [here](https://github.com/AlexGustafsson/systembolaget-api-fetch/releases/).

# Quickstart
<a name="quickstart"></a>

## Creating an account at systembolaget

Note: If you don't want to create an account, you may use the file-based APIs via the legacy branch or releases for versions lower than v2. The API may be deprecated by Systembolaget in the future without further notice. The legacy releases are no longer maintained but the file-based features may be added to future versions.

1. Go to https://api-portal.systembolaget.se and sign up
2. Subscribe to the open API via https://api-portal.systembolaget.se/products/Open%20API/subscribe
3. Access your access tokens via your profile https://api-portal.systembolaget.se/developer

## Using as a utility

Start by grabbing the latest release for your platform from the releases.

Download the assortment and print it as an XML to STDOUT.

```shell
./systembolaget download assortment --format=xml
```

Download the assortment and save it as prettified JSON to a file.

```shell
./systembolaget download assortment --pretty --output=assortment.json
```

Convert the assortment back to XML.

```shell
./systembolaget convert assortment --input assortment.json --output=assortment.xml
```

## Using as a library

Add the necessary import.

```go
import (
	systembolaget "github.com/alexgustafsson/systembolaget-api/systembolaget/v1"
)
```

Download assortment.

```go
response, err := systembolaget.DownloadAssortment()
```

Convert to prettified XML.

```go
xml, err := response.ConvertToXML(true)
```

Convert to minimal JSON.

```go
json, err := response.ConvertToJSON(false)
```

Load from JSON

```go
assortment, err := systembolaget.ParseAssortmentFromJSON(fileBytes)
```

# Table of contents

[Quickstart](#quickstart)<br/>
[Use cases](#usecases)<br/>
[Contributing](#contributing)<br/>
[Disclaimer](#disclaimer)

# Use cases
<a name="usecases"></a>

The utility can be used to automatically grab the latest available data in XML or JSON format. The data can be used to create interesting statistical charts, recommendation apps, search engines and more.

# Contributing
<a name="contributing"></a>

Any help with the project is more than welcome. This is my first take on Go and its ecosystem so some things might not be following best practices or be incorrectly implemented all together.

# Disclaimer
<a name="disclaimer"></a>

_Although the project is very capable, it is not built with production in mind. Therefore there might be complications when trying to use systembolaget-api for large-scale projects meant for the public. The utility was created to easily integrate Systembolaget's open APIs on several platforms and as such it might not promote best practices nor be performant._

_The author nor the utility is in any way affiliated with Systembolaget._
