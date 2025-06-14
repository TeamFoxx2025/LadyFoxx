# LadyFoxx Edge docs

Welcome to the LadyFoxx Edge documentation, built with [the Material theme for MkDocs](https://squidfunk.github.io/mkdocs-material/).

## Build and serve the site

### Clone the repo

```sh
https://github.com/TeamFoxx2025/LadyFoxx.git
cd foxx-chain/docs
```

### Run with Python

1. Download and install Python 3.11: https://www.python.org/downloads/

2. Install the `virtualenv` package:

```sh
pip install virtualenv
```

3. Build and serve the html

```sh
./run.sh
```

The site runs at: http://127.0.0.1:8000/

### Run with Docker

:warning: Remove line 10 from the `Dockerfile` to run locally.

1. Spin up the image:

```sh
docker build -t foxx-chain-docs .
```

2. Run the container:

```
docker compose up
```

The site runs at: http://127.0.0.1:8000/