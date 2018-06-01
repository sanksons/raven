## Raven

Does what Ravens are meant to do i.e deliver messages from one place to another.

Supports following engines:
- Redis
- Redis Cluster

## How to Install:

Simply run, below command
```
go get -u github.com/sanksons/raven
```

## How to use:

Detailed examples are kept in examples directory. But for quick view:

### Defining a Publisher:

Initialize Raven Farm.

```go
//
// Initialize raven farm.
//
farm, _ := raven.InitializeFarm(
    raven.FARM_TYPE_REDISCLUSTER,
    raven.RedisClusterConfig{
        Addrs:    []string{"172.17.0.2:30001"},
        PoolSize: 10,
    },
    nil,
)
```

Pick a raven from Farm.

```go
//Pick a raven from Farm
myraven := farm.GetRaven()
```

Hand over message to raven.

```go
//Hand over message to raven.
myraven.HandMessage(
    raven.PrepareMessage("msgID", "msgType","Message data!!"),
)
```

Specify destinationn for your raven.

```go
const DESTINATION = "product1"
const BUCKET = "1"
//Specify destinationn for your raven.
myraven.SetDestination(raven.CreateDestination(DESTINATION, BUCKET))
```

Make the Raven fly.

```go
//Make the Raven fly.
myraven.Fly()
```