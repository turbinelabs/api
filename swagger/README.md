# Swagger and Dredd

Swagger (http://swagger.io/) is a bunch of things, but at its core it's a way to specify a web api. Once you have this specification there are several tools to generate documentation, run validation tests, generate code, etc. It's also technically the Open API Specification (https://github.com/OAI/OpenAPI-Specification), but nobody apparently calls it that.

Dredd (https://dredd.readthedocs.io) is a tool that uses Swagger or API Blueprint (a competing spec format) to drive validation tests. Since this almost never works out of the box they have a way to define "hooks" that run before/after each case to massage stuff into a working form.

# Getting Started

The dredd getting started guide (https://dredd.readthedocs.io/en/latest/quickstart/) is pretty ok. Bascially just `npm install -g dredd` to get it on your path.

You'll also need /var/log/api-server created and write accessible for your user. If we log api-server to stdout (the default in the dev environment) it clutters dredd logs considerably.

With that done you should be able to run dredd in this directory (tbn/api/swagger) to run through a series of tests.

# A Bit More Depth

Configuration for dredd itself is handled via the dredd.yml file. What's configured right now is

* load hooks from the hooks.js file
* start up the api-server with `api-server --log.stdout=false`
* wait 1 second for server startup
* add the header "X-Turbine-API-Key: key-present" to every request
* run tests in sorted order (POST/GET/PUT/DELETE)
* set the blueprint file to swagger.yml
* set the endpoint under test to http://localhost:8080

Note that command line flags DO NOT override the dredd.yml settings. Which is annoying.
