# Velociraptor end to end test repo

This repository maintains an end to end test suite to ensure the
following end to end processes are working while migrating
Velociraptor release versions:

1. Ensure that data is still accessible when migrating from and older
   version of Velociraptor to a newer version. This also tests that
   basic operations (like searching for clients) work with the
   existing data.

2. Ensure that older clients interoperate with the newer server.
