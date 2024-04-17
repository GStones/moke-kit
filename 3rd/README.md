# Third Party

## [Agones](https://agones.dev/site/):

Host, Run and Scale dedicated game servers on Kubernetes

### Modules:
* `AgonesSDKModule` : agones go sdk module [(see more)](https://agones.dev/site/docs/guides/client-sdks/)
* `AgonesAllocateClientModule`: agones allocate grpc client module

### Environment Variables:
| ENV                     | Description                                                                               | Default                  |
|-------------------------|-------------------------------------------------------------------------------------------|--------------------------|
| AGONES_DEPLOYMENT       | agones deployment (local/dev/prod)<br/> local/dev: will mock a url with MOCK_ALLOCATE_URL | local                    |
| MOCK_ALLOCATE_URL       | mock allocate url(only for non-prod deployment)                                           | localhost:8888           |
| ALLOCATE_SERVICE_URL    | allocate service url(only for prod deployment)                                            | ""                       |
| ALLOCATE_CLIENT_CERT    | allocate client cert path                                                                 | ./configs/agones/tls.crt |
| ALLOCATE_CLIENT_KEY     | allocate client key path                                                                  | ./configs/agones/tls.key |
| ALLOCATE_SERVER_CA_CERT | allocate server ca cert path                                                              | ./configs/agones/ca.crt  |

## [IAP](https://github.com/awa/go-iap):

Verifies the purchase receipt via AppStore, GooglePlayStore or Amazon AppStore.

### Modules:
* `IAPModule`: iap module

### Environment Variables:
| ENV                    | Description                 | Default |
|------------------------|-----------------------------|---------|
| APPLE_KEY_ID           | apple key id                | ""      |
| APPLE_ISSUER           | apple issuer                | ""      |
| APPLE_PRIVATE_KEY      | apple private key path      | ""      |
| APPLE_BUNDLE_ID        | apple bundle id             | ""      |
| APPLE_SANDBOX          | apple sandbox               | true    |
| GOOGLE_PLAY_PUBLIC_KEY | google play public key path | ""      |

