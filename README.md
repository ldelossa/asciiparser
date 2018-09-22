# asciiparser

## WHAT
asciiparser takes an ascii body via HTTP post and accomplishes the following:
* report the total word count
* the count of each word occurence 
* report the size of the body
* creates a unique ID that can be used for retrievel
* removes any occurence of a word prefixed with "blue"

## HOW
The core components of this application are Go's standard HTTP server, an in memory store for holding upload records, and bufio library for stream processing the body.

## Getting it running
### Local Build
If you'd like to build the container locally you will first need to obtain dep for depdency management. You can find dep here https://github.com/golang/dep. 
Once installed you can run `dep ensure` in the root of the directory which will create a `vendor/` directory holding all dependencies.

With the vendor folder now present you are able to build the container locally with the following command:
```
docker build -f docker/Dockerfile -t ldelossa/asciiparser:v1.0.0 .
```
Feel free to define any repo/container:version string as you'd like. Once built you can run the container like so:
```
docker run -p 8080:8080 ldelossa/asciiparser:v1.0.0
```
again changing the repo/container:version if you used a different one during the build command.

You are now able to send post requests to "http://localhost:8080/api/v1/uploads"

### Docker Hub
You can also skip the local building above by pulling the container directly from my docker repository:

```
docker pull ldelossa/asciiparser
```

You will then use the same command to run the container
```
docker run -p 8080:8080 ldelossa/asciiparser:v1.0.0
```

## Usage
This service exposes an endpoint at `/api/v1/uploads`. Endpoint supports `POST` and `GET`. `GET` requests can be at the root of the path or you may add an ID to the endpoint. 

### Examples
POST
```
❯ curl -s -XPOST "http://localhost:8080/api/v1/uploads" --data "This is a really cool string blue" | jq .
{
  "size": 33,
  "id": "360d37d9-b3a7-441c-97f6-8f85d1a2b8da",
  "word_count": 7,
  "occurences": {
    "This": 1,
    "a": 1,
    "cool": 1,
    "is": 1,
    "really": 1,
    "string": 1
  }
  }
```

GET BY ID
```
❯ curl -s -XGET "http://localhost:8080/api/v1/uploads/360d37d9-b3a7-441c-97f6-8f85d1a2b8da" | jq .
{
  "size": 33,
  "id": "360d37d9-b3a7-441c-97f6-8f85d1a2b8da",
  "word_count": 7,
  "occurences": {
    "This": 1,
    "a": 1,
    "cool": 1,
    "is": 1,
    "really": 1,
    "string": 1
  }
}
```

GET ALL
```
❯ curl -s -XGET "http://localhost:8080/api/v1/uploads" | jq .
[
  {
    "size": 33,
    "id": "c9754287-5a0a-473d-9984-153f6470af85",
    "word_count": 7,
    "occurences": {
      "This": 1,
      "a": 1,
      "cool": 1,
      "is": 1,
      "really": 1,
      "string": 1
    }
  },
  {
    "size": 33,
    "id": "360d37d9-b3a7-441c-97f6-8f85d1a2b8da",
    "word_count": 7,
    "occurences": {
      "This": 1,
      "a": 1,
      "cool": 1,
      "is": 1,
      "really": 1,
      "string": 1
    }
  }
]
```
