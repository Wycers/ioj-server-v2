# Server-v2

<p align="center">
  <p align="center">
    <a href="https://github.com/infinity-oj/server-v2/releases/latest">
      <img alt="GitHub release" src="https://img.shields.io/github/v/release/infinity-oj/server-v2.svg?logo=github&style=for-the-badge" />
    </a>
    <a href="https://github.com/infinity-oj/server-v2/actions/workflows/build.yml">
       <img alt="Build workflow" src="https://img.shields.io/github/workflow/status/infinity-oj/server-v2/build?logo=github&style=for-the-badge" />
    </a>
  </p>
</p>

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
