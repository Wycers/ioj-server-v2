# Server-v2
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
