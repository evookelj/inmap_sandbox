# INMAP Sandbox

This repository uses evookelj/inmap to perform calculations about contribution and exposure to pollution by various income deciles and ethnicities.

## To Use
1. Properly set the variables in *setup.sh* (which should be self-explanatory).
2. ```cd ${INMAP_SANDBOX_ROOT}```
2. ```source setup.sh```
3. ```go run .```

Note that currently, this only runs on test data and exists for the purpose of creating the functionality and flow (rather than getting correct results). This is because a local machine does not have the capability to run the full model with the full data volume. The results of the original article were calculated using a "2018-vintage Google Compute Engine instance with 32 CPU cores, 208 GB of RAM, and a 500-GB hard drive."

## Files
- *data/*: holds various data files and configs necessary for running the sandbox.
- *contribution.go* provides functionality for calculating the pollution contribution of particular demographics
- *exposure.go* provides functionality for calculating the exposure to pollution of particular demographics
- *go.mod, go.sum* are standard files necessary for any Go module
- *main.go* calculates pollution exposure and contribution (using the functionality provided by the other files). This is where all running code should go
- *setup.sh* defines some environment variables necessary for proper functionality. Properly set these variables and source this script before running.
- *util.go* provides some utilities that aren't specific to exposure or contribution calculations

