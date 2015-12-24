# Gaia

[![Circle
CI](https://circleci.com/gh/NotyIm/gaia.svg?style=svg)](https://circleci.com/gh/NotyIm/gaia)

Gaia is NotyIM heart, it requests website, store data into InfluxDB.
Once it has data, it checks the data to see if it matches a set of
criteria. If yes, it creates an incident in system.

# Gaia interface

We interact with GaiA via a HTTP interface:


- POST: /monitor/
    - id:
    - address:

- DELETE: /monitor/{id}
    stop monitoring this website

- PUT: /monitor/
    - id:
    - address:
    Update

- 
