# Ethereum Swarm Chunk Size Distribution

This application prints distribution of Ethereum chunk sizes from Swarm's leveldb localstore retrieval data index.

It is compatible with github.com/ethersphere/bee v1.4.3.

The iteration is done over all chunks, so it may take a while in case of a larger database.

```sh
swarm-chunk-sizes -path /path/to/swar/data/dir/localstore
```
