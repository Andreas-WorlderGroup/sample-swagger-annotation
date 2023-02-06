
# Sample Project Golang Echo

This project contains 2 main file. The first one is core api which have api for get and store data. The second one is worker logic which have responsible to hit api store data with random value every x second.

## Tech Stack

**Backend:** Golang, Echo Framework, Mysql

## Deployment

For running core api
```bash
  make run-container-core
```

For running worker
```bash
  make run-worker
```

## Test api
Get api using ID
```bash
  curl --location --request GET 'localhost:8080/data?id1=4&id2=q'
```

Get api using timestamp
```bash
  curl --location --request GET 'localhost:8080/data?start_timestamp=1136214245&end_timestamp=1719672116'
```

## Badges

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](https://opensource.org/licenses/)
[![AGPL License](https://img.shields.io/badge/license-AGPL-blue.svg)](http://www.gnu.org/licenses/agpl-3.0)
