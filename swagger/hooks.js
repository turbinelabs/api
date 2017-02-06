var hooks = require('hooks');
var http = require('http');

var apiKey = "key-present"
var responseStash = {}

var defaultZoneName = "dredd-zone-" + Math.random().toString(36).substring(7);
var defaultDomainName = "dredd-domain-" + Math.random().toString(36).substring(7);
var defaultClusterName = "dredd-cluster-" + Math.random().toString(36).substring(7);
var defaultSharedRulesName = "dredd-shared-rules-" + Math.random().toString(36).substring(7);
var proxyName = "dredd-proxy-" + Math.random().toString(36).substring(7);
var routeName = "dredd-route-" + Math.random().toString(36).substring(7);
var clusterName = "dredd-cluster-" + Math.random().toString(36).substring(7);

// set up default objects. The default objects are used in cases where a test case (e.g. route) needs
// a pre-existng object to attach to. Other test cases (e.g. zone itself) are ordered, so they'll create,
// update, delete a specific domain for their test suite
hooks.beforeAll(function(transactions, done) {
  getOrMakeObject(transactions[0], "zone", defaultZoneName, {name: defaultZoneName}, function(response, body, err) {
    var zone = responseStash[defaultZoneName];
    var domain = {name: defaultDomainName, zone_key: zone.zone_key, port: 80}
    getOrMakeObject(transactions[0], "domain", defaultDomainName, domain, function(response, body, err) {
      var newDomain = responseStash[defaultDomainName];
      var cluster = {name: defaultClusterName, zone_key: zone.zone_key};
      getOrMakeObject(transactions[0], "cluster", defaultClusterName, cluster, function(response, body, err) {
        var newCluster = responseStash[defaultClusterName];
        var sr = {
          zone_key: zone.zone_key,
          default: {
            light: [{
              weight: 1,
              cluster_key: newCluster.cluster_key,
              metadata: [{key: "app", value: "hello-node"}],
              weight: 100
            }],
            dark: [{
              weight: 1,
              cluster_key: newCluster.cluster_key,
              metadata: [{key: "app", value: "hello-node"}],
              weight: 100
            }],
            tap: [{
              weight: 1,
              cluster_key: newCluster.cluster_key,
              metadata: [{key: "app", value: "hello-node"}],
              weight: 100
            }]
          },
          rules: [{
            rule_key: "foo",
            methods: ["GET"],
            matches: [{
              kind: "query",
              from: {key: "color", value: "blue"},
              to: {key: "color", value: "blue"}
            }],
            constraints: {
              light: [{
                weight: 1,
                cluster_key: newCluster.cluster_key,
                metadata: [{key: "app", value: "hello-node"}],
                weight: 100
              }],
              dark: [{
                weight: 1,
                cluster_key: newCluster.cluster_key,
                metadata: [{key: "app", value: "hello-node"}],
                weight: 100
              }],
              tap: [{
                weight: 1,
                cluster_key: newCluster.cluster_key,
                metadata: [{key: "app", value: "hello-node"}],
                weight: 100
              }]
            }
          }]
        }
        getOrMakeObject(transactions[0], "shared_rules", defaultSharedRulesName, sr, function(response, body, err) {
          done();
        });
      });
    });
  });
});

// zone hooks

var zoneName = "dredd-zone-" + Math.random().toString(36).substring(7);

// create a new zone, using a randomly generated name
hooks.before("Zone > /v1.0/zone > create zone > 200 > application/json", function(transaction) {
  var body = JSON.stringify({name: zoneName});
  transaction.request.body = body;
});

// in the after hook we parse out the response and stash it for later operations
hooks.after("Zone > /v1.0/zone > create zone > 200 > application/json", function(transaction) {
  zone = JSON.parse(transaction.real.body).result;
  responseStash[zoneName] = zone;
});

// mutate path to get the zone we created previously
hooks.before("Zone > /v1.0/zone/{zoneKey} > get zone > 200 > application/json", function(transaction) {
  dz = responseStash[zoneName];
  transaction.fullPath = "/v1.0/zone/" + dz.zone_key
  transaction.request.uri = "/v1.0/zone/" + dz.zone_key
});

// mutate the path to modify the zone we created previously
hooks.before("Zone > /v1.0/zone/{zoneKey} > modify zone > 200 > application/json", function(transaction) {
  var body = responseStash[zoneName];
  body.name = body.name + "-modified";
  transaction.request.body = JSON.stringify(body);
  transaction.fullPath = "/v1.0/zone/" + body.zone_key
  transaction.request.uri = "/v1.0/zone/" + body.zone_key
});

// in the after hook we stash the modified object so we can use the new checksum in the delete call
hooks.after("Zone > /v1.0/zone/{zoneKey} > modify zone > 200 > application/json", function(transaction) {
  zone = JSON.parse(transaction.real.body).result;
  responseStash[zoneName] = zone;
});

// mutate path to delete the zone we created previously
hooks.before("Zone > /v1.0/zone/{zoneKey} > delete zone > 200 > application/json", function(transaction) {
  ddz = responseStash[zoneName];
  transaction.fullPath = "/v1.0/zone/" + ddz.zone_key + "?checksum=" + ddz.checksum;
  transaction.request.uri = "/v1.0/zone/" + ddz.zone_key + "?checksum=" + ddz.checksum;;
});


// domain hooks
// these follow the zone hooks very closely, see comments there for an explanation of various hooks

var domainName = "dredd-domain-" + Math.random().toString(36).substring(7);

hooks.before("Domain > /v1.0/domain > create domain > 200 > application/json", function(transaction) {
  var zone = responseStash[defaultZoneName];
  var body = JSON.stringify({name: domainName, zone_key: zone.zone_key, port: 80});
  transaction.request.body = body;
});

hooks.after("Domain > /v1.0/domain > create domain > 200 > application/json", function(transaction) {
  domain = JSON.parse(transaction.real.body).result;
  responseStash[domainName] = domain;
});

hooks.before("Domain > /v1.0/domain/{domainKey} > get domain > 200 > application/json", function(transaction) {
  domain = responseStash[domainName];
  transaction.fullPath = "/v1.0/domain/" + domain.domain_key;
  transaction.request.uri = "/v1.0/domain/" + domain.domain_key;
});

hooks.before("Domain > /v1.0/domain/{domainKey} > modify domain > 200 > application/json", function(transaction) {
  var body = responseStash[domainName];
  body.name = body.name + "-modified";
  transaction.request.body = JSON.stringify(body);
  transaction.fullPath = "/v1.0/domain/" + body.domain_key
  transaction.request.uri = "/v1.0/domain/" + body.domain_key
});

hooks.after("Domain > /v1.0/domain/{domainKey} > modify domain > 200 > application/json", function(transaction) {
  domain = JSON.parse(transaction.real.body).result;
  responseStash[domainName] = domain;
});

hooks.before("Domain > /v1.0/domain/{domainKey} > delete domain > 200 > application/json", function(transaction) {
  domain = responseStash[domainName];
  transaction.fullPath = "/v1.0/domain/" + domain.domain_key + "?checksum=" + domain.checksum;
  transaction.request.uri = "/v1.0/domain/" + domain.domain_key + "?checksum=" + domain.checksum;;
});


// proxy hooks
// these follow the zone hooks very closely, see comments there for an explanation of various hooks

hooks.before("Proxy > /v1.0/proxy > create proxy > 200 > application/json", function(transaction) {
  var zone = responseStash[defaultZoneName];
  var domain = responseStash[defaultDomainName];
  var body = JSON.stringify({name: proxyName, zone_key: zone.zone_key, host: "foo.com", port: 80, domain_keys: [domain.domain_key]});
  transaction.request.body = body;
});

hooks.after("Proxy > /v1.0/proxy > create proxy > 200 > application/json", function(transaction) {
  proxy = JSON.parse(transaction.real.body).result;
  responseStash[proxyName] = proxy;
});

hooks.before("Proxy > /v1.0/proxy/{proxyKey} > get proxy > 200 > application/json", function(transaction) {
  proxy = responseStash[proxyName];
  transaction.fullPath = "/v1.0/proxy/" + proxy.proxy_key;
  transaction.request.uri = "/v1.0/proxy/" + proxy.proxy_key;
});

hooks.before("Proxy > /v1.0/proxy/{proxyKey} > modify proxy > 200 > application/json", function(transaction) {
  var body = responseStash[proxyName];
  body.name = body.name + "-modified";
  transaction.request.body = JSON.stringify(body);
  transaction.fullPath = "/v1.0/proxy/" + body.proxy_key
  transaction.request.uri = "/v1.0/proxy/" + body.proxy_key
});

hooks.after("Proxy > /v1.0/proxy/{proxyKey} > modify proxy > 200 > application/json", function(transaction) {
  proxy = JSON.parse(transaction.real.body).result;
  responseStash[proxyName] = proxy;
});

hooks.before("Proxy > /v1.0/proxy/{proxyKey} > delete proxy > 200 > application/json", function(transaction) {
  proxy = responseStash[proxyName];
  transaction.fullPath = "/v1.0/proxy/" + proxy.proxy_key + "?checksum=" + proxy.checksum;
  transaction.request.uri = "/v1.0/proxy/" + proxy.proxy_key + "?checksum=" + proxy.checksum;;
});

// shared rules hooks

var sharedRulesName = "dredd-shared-rules-" + Math.random().toString(36).substring(7);

// this hook is the exception. WHen we do a get, we're likely going to be pulling other fixture clusters that
// have empty elements. This fixes up any default members that are null and swaps them out with empty arrays
hooks.beforeValidation("Shared Rules > /v1.0/shared_rules > get shared_rules > 200 > application/json", function(transaction) {
  srs = JSON.parse(transaction.real.body);
  for (var i = 0; i < srs.result.length; i++) {
    var sr = srs.result[i];
    if (sr.default == null) {
      sr.default = {};
    } else {
      if (sr.default.dark == null) {
        sr.default.dark = [];
      }
      if (sr.default.tap == null) {
        sr.default.tap = [];
      }
      if (sr.default.light == null) {
        sr.default.light = [];
      }
    }
    if (sr.rules == null) {
      sr.rules = [];
    }
  }
  transaction.real.body = JSON.stringify(srs);
});

// create a new shared rules, using a randomly generated name
hooks.before("Shared Rules > /v1.0/shared_rules > create shared_rules > 200 > application/json", function(transaction) {
  var zone = responseStash[defaultZoneName];
  var cluster = responseStash[defaultClusterName];
  // this is verbose, but dredd really doesn't like empty arrays-as-nulls, so we provide values for everything.
  // it also results in better validation, as we can verify array elements
  var sr = {
    zone_key: zone.zone_key,
    default: {
      light: [{
        weight: 1,
        cluster_key: cluster.cluster_key,
        metadata: [{key: "app", value: "hello-node"}],
        weight: 100
      }],
      dark: [{
        weight: 1,
        cluster_key: cluster.cluster_key,
        metadata: [{key: "app", value: "hello-node"}],
        weight: 100
      }],
      tap: [{
        weight: 1,
        cluster_key: cluster.cluster_key,
        metadata: [{key: "app", value: "hello-node"}],
        weight: 100
      }]
    },
    rules: [{
      rule_key: "foo",
      methods: ["GET"],
      matches: [{
        kind: "query",
        from: {key: "color", value: "blue"},
        to: {key: "color", value: "blue"}
      }],
      constraints: {
        light: [{
          weight: 1,
          cluster_key: cluster.cluster_key,
          metadata: [{key: "app", value: "hello-node"}],
          weight: 100
        }],
        dark: [{
          weight: 1,
          cluster_key: cluster.cluster_key,
          metadata: [{key: "app", value: "hello-node"}],
          weight: 100
        }],
        tap: [{
          weight: 1,
          cluster_key: cluster.cluster_key,
          metadata: [{key: "app", value: "hello-node"}],
          weight: 100
        }]
      }
    }]
  }
  var body = JSON.stringify(sr);
  transaction.request.body = body;
});

// in the after hook we parse out the response and stash it for later operations
hooks.after("Shared Rules > /v1.0/shared_rules > create shared_rules > 200 > application/json", function(transaction) {
  sharedRules = JSON.parse(transaction.real.body).result;
  responseStash[sharedRulesName] = sharedRules;
});

// mutate path to get the shared rules we created previously
hooks.before("Shared Rules > /v1.0/shared_rules/{sharedRulesKey} > get shared_rules object > 200 > application/json", function(transaction) {
  sr = responseStash[sharedRulesName];
  transaction.fullPath = "/v1.0/shared_rules/" + sr.shared_rules_key
  transaction.request.uri = "/v1.0/shared_rules/" + sr.shared_rules_key
});

// mutate the path to modify the shared rules we created previously
hooks.before("Shared Rules > /v1.0/shared_rules/{sharedRulesKey} > modify shared_rules object > 200 > application/json", function(transaction) {
  var body = responseStash[sharedRulesName];
  transaction.request.body = JSON.stringify(body);
  transaction.fullPath = "/v1.0/shared_rules/" + body.shared_rules_key
  transaction.request.uri = "/v1.0/shared_rules/" + body.shared_rules_key
});

// in the after hook we stash the modified object so we can use the new checksum in the delete call
hooks.after("Shared Rules > /v1.0/shared_rules/{sharedRulesKey} > modify shared_rules object > 200 > application/json", function(transaction) {
  sr = JSON.parse(transaction.real.body).result;
  responseStash[sharedRulesName] = sr;
});

// mutate path to delete the zone we created previously
hooks.before("Shared Rules > /v1.0/shared_rules/{sharedRulesKey} > delete shared_rules object > 200 > application/json", function(transaction) {
  dsr = responseStash[sharedRulesName];
  transaction.fullPath = "/v1.0/shared_rules/" + dsr.shared_rules_key + "?checksum=" + dsr.checksum;
  transaction.request.uri = "/v1.0/shared_rules/" + dsr.shared_rules_key + "?checksum=" + dsr.checksum;;
});

// route hooks
// these follow the zone hooks very closely, see comments there for an explanation of various hooks

// this hook is the exception. WHen we do a get, we're likely going to be pulling other fixture routes that
// have empty elements. This fixes up any default and rules members that are null and swaps them out with empties
hooks.beforeValidation("Route > /v1.0/route > get routes > 200 > application/json", function(transaction) {
  routes = JSON.parse(transaction.real.body);
  for (var i = 0; i < routes.result.length; i++) {
    route = routes.result[i];
    if (route.rules == null) {
      route.rules = [];
    }
  }
  transaction.real.body = JSON.stringify(routes);
});

hooks.before("Route > /v1.0/route > create route > 200 > application/json", function(transaction) {
  var zone = responseStash[defaultZoneName];
  var domain = responseStash[defaultDomainName];
  var cluster = responseStash[defaultClusterName];
  var sr = responseStash[defaultSharedRulesName];
  // this is verbose, but dredd really doesn't like empty arrays-as-nulls, so we provide values for everything.
  // it also results in better validation, as we can verify array elements
  var body = JSON.stringify({name: routeName,
                             zone_key: zone.zone_key,
                             domain_key: domain.domain_key,
                             shared_rules_key: sr.shared_rules_key,
                             path: "/",
                             rules: [{
                               rule_key: "foo",
                               methods: ["GET"],
                               matches: [{
                                 kind: "query",
                                 from: {key: "color", value: "blue"},
                                 to: {key: "color", value: "blue"}
                               }],
                               constraints: {
                                 light: [{
                                   weight: 1,
                                   cluster_key: cluster.cluster_key,
                                   metadata: [{key: "app", value: "hello-node"}],
                                   weight: 100
                                 }],
                                 dark: [{
                                   weight: 1,
                                   cluster_key: cluster.cluster_key,
                                   metadata: [{key: "app", value: "hello-node"}],
                                   weight: 100
                                 }],
                                 tap: [{
                                   weight: 1,
                                   cluster_key: cluster.cluster_key,
                                   metadata: [{key: "app", value: "hello-node"}],
                                   weight: 100
                                 }]
                               }
                             }]
                            });
  transaction.request.body = body;
});

hooks.after("Route > /v1.0/route > create route > 200 > application/json", function(transaction, done) {
  // refresh cluster def in the responseStash because creating a route bumps the checksum
  getObject(transaction, "cluster", defaultClusterName, function(response, body, err) {
    done();
  });
  route = JSON.parse(transaction.real.body).result;
  responseStash[routeName] = route;
});

hooks.before("Route > /v1.0/route/{routeKey} > get route > 200 > application/json", function(transaction) {
  route = responseStash[routeName];
  transaction.fullPath = "/v1.0/route/" + route.route_key;
  transaction.request.uri = "/v1.0/route/" + route.route_key;
});

hooks.before("Route > /v1.0/route/{routeKey} > modify route > 200 > application/json", function(transaction) {
  var body = responseStash[routeName];
  body.name = body.name + "-modified";
  transaction.request.body = JSON.stringify(body);
  transaction.fullPath = "/v1.0/route/" + body.route_key
  transaction.request.uri = "/v1.0/route/" + body.route_key
});

hooks.after("Route > /v1.0/route/{routeKey} > modify route > 200 > application/json", function(transaction) {
  route = JSON.parse(transaction.real.body).result;
  responseStash[routeName] = route;
});

hooks.before("Route > /v1.0/route/{routeKey} > delete route > 200 > application/json", function(transaction) {
  route = responseStash[routeName];
  transaction.fullPath = "/v1.0/route/" + route.route_key + "?checksum=" + route.checksum;
  transaction.request.uri = "/v1.0/route/" + route.route_key + "?checksum=" + route.checksum;;
});

// cluster hooks
// most of these follow the zone hooks very closely, see comments there for an explanation of various hooks


// this hook is the exception. WHen we do a get, we're likely going to be pulling other fixture clusters that
// have empty elements. This fixes up any instances members that are null and swaps them out with empty arrays
hooks.beforeValidation("Cluster > /v1.0/cluster > get clusters > 200 > application/json", function(transaction) {
  clusters = JSON.parse(transaction.real.body);
  for (var i = 0; i < clusters.result.length; i++) {
    cluster = clusters.result[i];
    if (cluster.instances == null) {
      cluster.instances = [];
    }
  }
  transaction.real.body = JSON.stringify(clusters);
});

hooks.before("Cluster > /v1.0/cluster > create cluster > 200 > application/json", function(transaction) {
  var zone = responseStash[defaultZoneName];
  var body = JSON.stringify({name: clusterName,
                             zone_key: zone.zone_key,
                             instances: [
                               {host: "foo.bar.com",
                                port: 80,
                                metadata: [
                                  {key: "color", value: "blue"}
                                ]
                               }
                             ]});
  transaction.request.body = body;
});

hooks.after("Cluster > /v1.0/cluster > create cluster > 200 > application/json", function(transaction) {
  cluster = JSON.parse(transaction.real.body).result;
  responseStash[clusterName] = cluster;
});

hooks.before("Cluster > /v1.0/cluster/{clusterKey} > get cluster > 200 > application/json", function(transaction) {
  cluster = responseStash[clusterName];
  transaction.fullPath = "/v1.0/cluster/" + cluster.cluster_key;
  transaction.request.uri = "/v1.0/cluster/" + cluster.cluster_key;
});

hooks.before("Cluster > /v1.0/cluster/{clusterKey} > modify cluster > 200 > application/json", function(transaction) {
  var body = responseStash[clusterName];
  body.name = body.name + "-modified";
  transaction.request.body = JSON.stringify(body);
  transaction.fullPath = "/v1.0/cluster/" + body.cluster_key
  transaction.request.uri = "/v1.0/cluster/" + body.cluster_key
});

hooks.after("Cluster > /v1.0/cluster/{clusterKey} > modify cluster > 200 > application/json", function(transaction) {
  cluster = JSON.parse(transaction.real.body).result;
  responseStash[clusterName] = cluster;
});

hooks.before("Cluster > /v1.0/cluster/{clusterKey} > delete cluster > 200 > application/json", function(transaction) {
  cluster = responseStash[clusterName];
  transaction.fullPath = "/v1.0/cluster/" + cluster.cluster_key + "?checksum=" + cluster.checksum;
  transaction.request.uri = "/v1.0/cluster/" + cluster.cluster_key + "?checksum=" + cluster.checksum;;
});

// cluster -> instances hooks

// similar to zone hooks post methods
hooks.before("Cluster > /v1.0/cluster/{clusterKey}/instances > add instance > 200 > application/json", function(transaction) {
  var cluster = responseStash[defaultClusterName];
  transaction.fullPath = "/v1.0/cluster/" + cluster.cluster_key + "/instance?checksum=" + cluster.checksum;
  transaction.request.uri = "/v1.0/cluster/" + cluster.cluster_key + "/instance=" + cluster.checksum;
  var body = {
    host: "bar1.bar.com",
    port: 80,
    metadata: [{key: "color", value: "blue"}]
  }
  transaction.request.body = JSON.stringify(body);
});

// after we make this host we need to refresh the cluster definition to pull int he new instances and update th checksum
hooks.after("Cluster > /v1.0/cluster/{clusterKey}/instances > add instance > 200 > application/json", function(transaction, done) {
  // refresh cluster def in the responseStash because adding an instance bumps the checksum
  getObject(transaction, "cluster", defaultClusterName, function(response, body, err) {
    done();
  });
  responseStash["clusterInstance"] = JSON.parse(transaction.real.body).result.instances[0];
});

hooks.before("Cluster > /v1.0/cluster/{clusterKey}/instances/{instanceIdentifier} > remove instance > 200 > application/json", function(transaction) {
  var cluster = responseStash[defaultClusterName];
  var inst = responseStash["clusterInstance"];
  var path = "/v1.0/cluster/" + cluster.cluster_key + "/instance/" + inst.host + ":" + inst.port + "?checksum=" + cluster.checksum;
  transaction.fullPath = path;
  transaction.request.uri = path;
});


// utility functiosn

// look up an object of a given type from the general GET path (e.g. /v1.0/zones).
// if found, update responseStash with an entry at objectName. After completion
// call andThen(response, responseBody, err)
function getObject(transaction, objectType, objectName, andThen) {
  var options = {
    host: transaction.host,
    port: transaction.port,
    path: "/v1.0/" + objectType,
    method: "GET",
    headers: {
      "X-Turbine-API-Key": apiKey
    }
  }
  callback = function(response) {
    var str = '';
    if (response.statusCode == 200) {
      response.on('data', function(chunk) {
        str += chunk;
      });
      response.on('end', function() {
        try {
          var obj = JSON.parse(str);
          for (var i = 0; i < obj.result.length; i++) {
            if (obj.result[i].name == objectName) {
              var elem = obj.result[i];
              responseStash[objectName] = elem;
            }
          }
        } catch(err) {
          andThen(null, null, err);
        }
        if (andThen) {
          andThen(response, str, null);
        }
      });
    } else {
      if (andThen) {
        andThen(response, str, null);
      }
    }
  }
  var req = http.request(options, callback);
  req.write("");
  req.end();
}

// create an object of a given type from the supplied newObject definition.
// If successful create an entry in responseStash for the given objectName.
// After completion call andThen(response, responseBody, err)
function makeObject(transaction, objectType, objectName, newObject, andThen) {
  var options = {
    host: transaction.host,
    port: transaction.port,
    path: "/v1.0/" + objectType,
    method: "POST",
    headers: {
      "X-Turbine-API-Key": apiKey
    }
  }

  callback = function(response) {
    var str = '';
    response.on('data', function(chunk) {
      str += chunk;
    });
    response.on('end', function() {
      try {
        var obj = JSON.parse(str);
        responseStash[objectName] = obj.result;
      } catch(err) {
        andThen(null, null, err);
      }
      if (andThen) {
        andThen(response, str, null);
      }
    });
  }
  var req = http.request(options, callback);
  req.write(JSON.stringify(newObject));
  req.end();
}

// helper to chain a get and a make call, creating entries that don't exist.
// in practice we shouldn't collide with old definitions.
function getOrMakeObject(transaction, objectType, objectName, newObject, andThen) {
  getObject(transaction, objectType, objectName, function(response, body, err) {
    if (responseStash[objectName]) {
      if (andThen) {
        andThen(response, body, null);
      }
    } else {
      makeObject(transaction, objectType, objectName, newObject, andThen);
    }
  });
}
