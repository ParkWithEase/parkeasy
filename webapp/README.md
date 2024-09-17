## ParkEasy Web Edition

The browser frontend for ParkEasy

### Getting Started

First, install all dependencies

    npm install

Then, a development server can be started with

    npm run dev

### Tests

Tests are divided into integration and unit tests

    npm run test:unit        # Run unit tests
    npm run test:integration # Run integration tests

    npm run test # Run all tests

### Linting

Linting rules are tested by CI and will block merges should they fail.

Linters can be run manually via

    npm run lint

Code formatting consistency is applied by [Prettier] and can be run manually via

    npm run format

It is recommended to use [integrations](https://prettier.io/docs/en/editors.html) available for your editor to remove the need to run this step.

[Prettier]: https://prettier.io/
