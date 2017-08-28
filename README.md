systembolaget-api-fetch  - a cross-platform solution for fetching and converting [Systembolaget's open APIs](https://www.systembolaget.se/api/) to JSON
======

The utility is an easy to use way of fetching Systembolaget's open APIs and converting them to JSON. It's written in Go and is available via several cross-platform builds.

# Quickstart
<a name="quickstart"></a>

#### Using systembolaget-api-fetch

Start by grabbing the latest release for your platform from the releases.

```
# run systembolaget-api-fetch
> ./systembolaget-api-fetch
Downloading https://www.systembolaget.se/api/assortment/stores/xml
...
Successfully processed stores API
Downloading https://www.systembolaget.se/api/assortment/products/xml
...
Successfully processed assortment API
Downloading https://www.systembolaget.se/api/assortment/stock/xml
...
Successfully processed inventory API

# use generated data stored in ./output/
> tree output
output
├── json
│   ├── assortment.json
│   ├── inventory.json
│   └── stores.json
└── xml
    ├── assortment.xml
    ├── inventory.xml
    └── stores.xml
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

_Although the project is very capable, it is not built with production in mind. Therefore there might be complications when trying to use systembolaget-api-fetch for large-scale projects meant for the public. The utility was created to easily integrate Systembolaget's open APIs on several platforms and as such it might not promote best practices nor be performant._

_The author nor the utility is in any way affiliated with Systembolaget._
