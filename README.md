# BCache

A lightweight cache over http

## API

### Endpoint
`/v1/{bucket}/{key}`

### Methods

#### PUT/POST
`curl -XPUT -d "This is a test value" 127.0.0.1:3017/v1/cache1/deadb33f`

#### GET
`curl -XGET 127.0.0.1:3017/v1/cache1/deadb33f`

#### DELETE
`curl -XDELETE 127.0.0.1:3017/v1/cache1/deadb33f`

