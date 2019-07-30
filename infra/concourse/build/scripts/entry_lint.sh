#!/bin/bash
set -e

source test/ci_integration.sh
setup_environment
make check
