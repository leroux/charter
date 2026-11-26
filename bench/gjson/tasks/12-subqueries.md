kind: neutral
I would like to nest queries. With

    {"topology":{"instances":[
      {"service_version":"1.2.3","service_roles":["one","two"]},
      {"service_version":"1.2.4","service_roles":["three","four"]},
      {"service_version":"1.2.2","service_roles":["one"]}]}}

the path `topology.instances.#(service_roles.#(=="one"))#.service_version` should return `["1.2.3","1.2.2"]`: the inner query is the query path of the outer, and an outer query with no operator matches the elements where its path exists. Whitespace around the inner query, as in `#( service_roles.#(=="one"))#`, is allowed, and quoted values in queries may contain escaped quotes and parentheses. Please implement it, mention nested queries in the docs, and add tests.
