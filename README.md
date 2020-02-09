# Url Shortener

[![CircleCI](https://circleci.com/gh/dharnitski/url-shortener.svg?style=svg)](https://circleci.com/gh/dharnitski/url-shortener)
[![Coverage Status](https://coveralls.io/repos/github/dharnitski/url-shortener/badge.svg)](https://coveralls.io/github/dharnitski/url-shortener)

## Design Limitations

* Site provides Anonymous access to functionality. User security is not implemented.
* Site generates very small urls by default after setup. It is not considered as an issue according to requirements but can be easily fixed by small change in DB scaffolding script.
* Users can predict previous/next generated URL. Not considered as an issue as there is no security requirements to prevent that and prediction does not affect site performance.

## Architecture

* Different technologies for front-end and back-end to use right tools.
* front-end is hosted by back-end to simplify setup and solve same origin issue.
* Unit tests to cover app business logic.
* Integration tests to cover SQL queries against real DB.
* Solution is configured to run unit tests by default to simplify onboarding for new developers.
* multi setup - native for development to speedup code -> test cycles, docker for deployments, circle-ci native to speed-up CI cycles (ci images are loaded faster).

