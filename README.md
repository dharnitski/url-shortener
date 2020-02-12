# Url Shortener

[![CircleCI](https://circleci.com/gh/dharnitski/url-shortener.svg?style=svg)](https://circleci.com/gh/dharnitski/url-shortener)
[![Coverage Status](https://coveralls.io/repos/github/dharnitski/url-shortener/badge.svg)](https://coveralls.io/github/dharnitski/url-shortener)

## Prerequisites for running app

* Docker and Docker Compose
* Port 3306 should be open

## Prerequisites for Development and Running tests

See `/client` and `/server` folders for more details

## Design Limitations

* Site provides Anonymous access to functionality. User security is not implemented.
* Site generates very small urls by default after setup. It is not considered as an issue according to requirements and can be easily fixed by small change in DB scaffolding script.
* Users can predict previous/next generated URL. Not considered as an issue as there is no security requirements to prevent users to see someone else data.
* Implementation assumes that front-end compiled assets does not match shorten URL pattern

## Architecture

* Different technologies are used for front-end and back-end to address different problems with right tool.
* Solution implemented close to 'vanilla' Angular and Go to reduce learning time for new developers
* front-end is hosted by back-end to simplify setup and solve same origin issue.
* Unit tests added to cover app business logic.
* Integration tests addd to validate SQL queries against real DB.
* Solution is configured to run unit tests by default to simplify DEV environment setup for new developers.
* Multi setup - native for development, docker for deployments, circle-ci native for CI.
* Database scaffolded using SQL migration scripts embedded into application and integration tests.
* BackEnd uses MySQL for storing data. Storage can be replaced with different SQL or NoSQL implementation.
* ShortUrls are generated as base 62 encoded string for record identifier. Hashing of original URL is not used due to its length and potential collisions.
