# Single-Cache System

This repository contains multiple separate packages that, together, create a single-cache set up with a server (acting as the backing store) and a single cache, to which requests are sent.

Both cache and server defaultly use localhost on on ports 4040 and 4000, respectively.

## Cache Endpoints:
`/c-dump` -- dumps all the items in the cache.
`/request/{item}` -- request an item. If the item is not in the cache, the cache automatically forwards the request to the backing store. 

## Backingstore Endpoints
`/b-dump` -- dumps all items in B.
`/fetch/{item}` -- Like `/request/` but for B.

## Policies
As of 23 Sept 2026, only the LRU policy is supported. More policies can be added in `policies.go`.