# Url Shortener

[![CircleCI](https://circleci.com/gh/dharnitski/url-shortener.svg?style=svg)](https://circleci.com/gh/dharnitski/url-shortener)
[![Coverage Status](https://coveralls.io/repos/github/dharnitski/url-shortener/badge.svg)](https://coveralls.io/github/dharnitski/url-shortener)

## Design Limitations

* Site provides Anonymous access to functionality. User security is not implemented.
* Site generates very small urls by default after setup. It is not considered as an issue according to requirements but can be easily fixed by small change in DB scaffolding script.
* Users can predict next generated URL. Not considered as an issue as there is no security requirements to prevent that and prediction does not affect site performance.





