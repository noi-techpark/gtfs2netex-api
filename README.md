# gtfs2netex-api
Small API wrapping https://github.com/noi-techpark/GTFS2NeTEx-converter with a REST call.

# Usage
Pass in a GTFS file with some configuration parameters, and you get back a NeTEx compliant `.xml` file

The container has no configuration, just make a http POST call to `/` with the following `multipart/form-data` (see [calls.http](./calls.http)):

```
nuts (string)
vat (string)
version (string)
az (string)
file (GTFS file - zip)
```

All parameters are required, and are passed on to the corresponding converter tool arguments.  
For more details on the parameters, refer to the converter's documentation

# Build and run
Before being able to build or run the `docker-compose` locally, make sure to check you the submodule using `git submodule update --init --recursive`

If you run `main.go` directly, it requires one single command line argument, which is the path to the converter tool main file (e.g. `./GTFS2NeTEx-converter/GTFS2NeTEx-converter.py`)