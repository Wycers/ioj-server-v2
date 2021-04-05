# Server-v2

[![Build](https://github.com/infinity-oj/server-v2/actions/workflows/build.yml/badge.svg)](https://github.com/infinity-oj/server-v2/actions/workflows/build.yml)

Server for Infinity OJ

## Technology Stack
1. Web Framework: Gin

## Prerequisites
1. [Postgres](https://www.postgresql.org/)

## Development
1. Run postgres service.
   for the database:
``` postgresql
create type judge_status as enum ('Pending', 'PartiallyCorrect', 'WrongAnswer', 'Accepted', 'SystemError', 'JudgementFailed', 'CompilationError', 'FileError', 'RuntimeError', 'TimeLimitExceeded', 'MemoryLimitExceeded', 'OutputLimitExceeded', 'InvalidInteraction', 'ConfigurationError', 'Canceled');
```
## Production
``` bash
make prod
```

## Usage
