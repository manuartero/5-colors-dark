# Unica Argentina

<p align="center">
  <img src="./icon.png"/>
</p>

This repository contains the implementation of the
Argentina's UNICA interface; it also contains a `TService` for direct access to some third parties in Movistar Argentina which access won't be done through UNICA.


## :book: Service interfaces

<p align="center">
  <img src="./doc/service-interfaces.png"/>
</p>

### REST interfaces

Argentina implementation of the following Unica interfaces. See [UnicaGateway](https://doc.tuenti.io/repos/unica-gateway/).

- `UnicaBillingEvents`
- `UnicaConsumption`
- `UnicaCustomerAccounts`
- `UnicaInvoicing`
- `UnicaOauth`
- `UnicaPrepayBalance`
- `UnicaProductCatalog`
- `UnicaProductInventory`
- `UnicaProductOrder`
- `UnicaUserPermissions`
- ~~`UnicaUserInfo`~~ **_not implemented yet_**.

### TService interface

Functionalities very specific for Argentina which are not covered by the UNICA interface. Consumed by other Novum services different from the Gateway (for example [UnicaCharging](https://gitlab.tuenti.io/tuenti/unica-charging) or [AccountDashboard](https://gitlab.tuenti.io/account/dashboard)).

All of those functionalities are covered by the same service: [`MovistarArIntegrationService`](http://artifacts.tuenti.io/idl/MovistarArIntegrationService.v1/); main groups of functionality covered by this interface are:

* **Ask For Balance**: additional type of Balance Transfer (see `UnicaPrepayBalance`) where the requester asks another customer to give him balance.
* **Club Movistar**: external webview developed by a third party with special offers for Movistar customers.
* **Company Management**: b2b feature; implementation of the idea of roles and managers.
* **Company Management (Roaming)**: b2b feature; one of the specific functionalities for a manager.
* **Connection Speed**: specific service for getting if the user is connected to the network, and its bandwidth status.
* **Familia Movistar**: external webview developed by a third party.
* **IxD** (a.k.a Internet Por Dia): a product for prepay customers that provides a data quota (MBs). Although it's mainly included in the Unica's response, it's also included here.
* **Landline**: b2b feature; ownership verification of a landline.
* **Mercado Pago**: kind of PayPal for Argentina customers: a third party service for performing card payments.
* **Multiply Balance**: a product provided to prepay customers which multiply the balance the customer topups in a period of time, most of it is configured manually in our side.
* **Plan Change**: interaction for changing the plan of a customer. In particular, returns the suggested best plan by Movistar Argentina, request a plan change to the suggested plan and check of the status.
* **Sms Premium**: specific service
* **SVA**: acronym for "_Servicio de Valor Añadido_" specific packages of Movistar Argentina.
* **SVA (Auto-Topup)**: a service provided to prepay customers, able to check status and request status changes.
* **Test Drive**: the FlagShip-Feature of Movistar Argentina. A three month long promotion that offers the customer a considerable data quota (~10Gb) and call mins over the promotion duration. At the end of the promotional period, a plan is suggested to the user, expecting a plan change request.


## :calling: Third party dependencies

We have a [detailed documentation of the integration with the OB](./doc/ob-integration-reference.md); but if just look for a general idea of the project dependencies:

**Boss API**: Deprecated API of Movistar Argentina: only used to get the token

**ApiGateway**: REST API providing most of the functionalities related with B2C customers and assuming slowly all functionalities deprecated at **Boss Api**

    - `/customerAccounts/v2` & `/userinfo/t3`: user profiling
    - `/balancemanagement/v3`: current balance of a b2c user
    - `/facturacion/v4`: billing information
    - `/movics/v1`: holds different functionalities related with user accounts
        - `/balances`: returns the current balance of a b2b user
        - `/roaming`: roaming operations (enable and disable)
    - `/mercadopago`: third party related with card payment (similar to PayPal)
    - `/loadbalance` & `/sosrecharge` & `/recarga`: topup operations
    - `/prepaidoffer/v2`: api for recurrent packages (an special type of bundles for prepay customers with recurrent payment)
    - `/svagw/facturacion`: billing information
    - `/userCenter/v1`: apoderado credential check
    - `/svagw/userinfo`: user profiling
    - `/svagw/ocs`: information related with the user consumption;
        - `/svagw/ocs/movimientos` billing events
        - `/svagw/ocs/bonos` consumption counters
        - `/svagw/ocs/monederos` buckets of a prepay customer
    - `/svagw/productoffering`: tariff details
    - `/svagw/sva`: used for several services such as "AutoTopup" service, "Apn" messages... usually implies the user receiving an SMS. Apart from that, is the entry point for the Catalog of regular packages for prepay customers (ProductCatalog and ProductOrder).

**Fourth Platform API**: REST API that will replace in theory any other third party integration (...)

    - `/products`: [subscribed products](https://docs.baikalplatform.com/subscribed_products/index.html); used for retrieving all the account types of a user (landline, mobile, internet connection...)

**SMS Premium API**: REST API providing functionality to handle an specific functionallity related with Premium SMS.

**Sispro**: REST API providing functionality to handle an specific functionallity related with plan change.

**ClubMovistar**: REST API providing functionality to handle the specific functionallity "Club Movistar".

**SOAP services**: Bundle of SOAP services related with B2B's landline requests.

**Public websites**: The following features are reached with normal HTTP requests:

    - `movistar.com.ar`: used for retrieving information related with the current tariff of the user
    - `ofertaideal.com`: used for retrieving the ideal plan for a user.

**[Unica Gateway](https://gitlab.tuenti.io/tuenti/unica-gateway)**: (internal dependency) Depending of the request, Gateway will be the requester or a dependency from the point of view of `UnicaAr`. Check the [authentication section](#authentication) for details.

### Authentication

Currently we have these scenarios:

 1. **`UnicaGw`** requests any of the Unica services: the request will be already authenticated with a valid token from the third party.
 Here, `UnicaGw` would have done two requests in series; the first to `UnicaArgentina.UnicaOauth` in order to obtain a valid token for the second request.
 Note that the first request won't be done if a valid token is already cached by `UnicaGw`
 2. **Other service at Novum** requests to the `tService` interface: `UnicaAr` will call `UnicaGw` requesting the authentication token (this is due to the cache that has `UnicaGw` containing the auth tokens for all countries); if the token must be negotiated again, `UnicaGw` will call back `UnicaArgentina.UnicaOauth` implementation. `UnicaAr` resolves the original request.
 3. **`UnicaGw`** requests a Unica services that dependes on 4p: the request will be already authenticated with a valid token from 4p.
 This time, `UnicaGw` would have obtained the token from the `LoginService` (another service at NOVUM).

<p align="center">
<img src="./doc/authentication.png"/>
</p>


## :chart_with_upwards_trend: Operational information

### Operational monitoring

**Argentina B2B instance**

* [Third parties Grafana dashboard and Alerts](https://metrics.tuenti.io/dashboard/db/novum-argentina-b2b-third-parties)
* [Service Grafana dashboard](https://metrics.tuenti.io/dashboard/db/k8s-unica-ar-implementation-movistar-b2b-ar-prod-service)
* [Kibana dashboard](https://logger.tuenti.io/index.html#/dashboard/elasticsearch/BSS%20Template?serviceName=unica-ar-implementation)
* [K8s deployment info](https://fisheye.tuenti.io/browse/k8s-definitions/services/unica-implementation/mad.movistar-b2b-ar-prod.unica-ar-implementation.kubedeploy.yaml?hb=true)

**Argentina B2C instance**

* [Third parties Grafana dashboard and Alerts](https://metrics.tuenti.io/dashboard/db/novum-argentina-third-parties-miami)
* [Service Grafana dashboard](https://metrics.tuenti.io/dashboard/db/k8s-java-service-mia-movistar-ar?var-smoothing=1&var-service=unica-ar-implementation&var-serviceName=unica-ar-implementation&var-nginxService=unica_ar_implementation)
* [Kibana dashboard](https://kibana.prd-mia.tuenti.io/app/kibana#/discover/d9a239a0-97e2-11e7-a8e1-d98564e7d4ac)
* [K8s deployment info](https://fisheye.tuenti.io/browse/k8s-definitions/services/unica-implementation/novum-mia-prd.movistar-ar.unica-implementation.kubedeploy.yaml?hb=true)

### Operational procedures:

#### Database config

Config files are: [B2B](https://gitlab.tuenti.io/tuenti/unified_config/blob/master/environment/tuenti/base/unica-ar-implementation/databaseConfig.yaml) & [B2C](https://gitlab.tuenti.io/tuenti/unified_config/blob/master/novum/movistar-ar/unica-ar-implementation-b2c/databaseConfig.yaml)

#### Local storage

In general terms, this service guarantees that all data provided is a fresh copy requested to the Argentina API,
and will be responsability for the calling service to implement a cache strategy.

**Exceptions**: the following requests are actually cached by `UnicaArgentina`:

 * Tariff details
    - **Api Gateway Api**: response from `/svagw/productoffering`
    - **Json Plan Api**: response from `/plan`
 * Invoicing
   - **Api Gateway API**: response from `/svagw/facturacion/periodos`
 * SMS premium
   - **Oauth SMS client**: response from `/token.php`

#### Service Config

The config of the service is defined in [unified_config](https://git.code.tuenti.io/unified_config).
Paths are:

  - **B2C**: `environment/novum/movistar-ar[-next|-cert|-dev]/unica-ar-implementation-b2c`
  - **B2B**: `environment/tuenti/[base|pre]/unica-ar-implementation`


## :gear: Development

* Set up your development environment.
* Write your code in a new branch.
* Run tests and check your changes locally.
* Deploy to preproduction and verify your changes.
* Open a code review.
* Integrate your changes into master branch.

### Environment setup

Follow instructions for [Java environment installation & setup](https://doc.tuenti.io/global-doc/platform/java/setup/).

To have the IDE correctly support lombok annotations:
- In the menu, select **File** > **Settings**
- In the options tree, go to **Plugins**
- Search for **Lombok Plugin** and install it

### Running locally

We've created two guides for running this repo locally;
  - the first one uses mocks for external dependencies so you won't reach any real API. This is most secure option.
  - The second one will hit the real API while keeping your service in your machine which is useful for local debuging. </br>
    Note that Argentina doesn't have a full implemented dev environment, so you'll using the Argentina PRO endpoints: be aware of non dry operations as purchases.

#### Mock option: everything keeps local

<p align="center">
<img src="./doc/running-locally.png"/>
</p>

**1. Mocking UnicaOAuth service** (Unica Gateway): the script `/scripts/tservicestub.groovy` launchs in your local machine a server at `localhost:5050` with a mocked UnicaOauth implementation.

```bash
$> groovy scripts/local-dev/tServiceStub.groovy
  ...
  INFORMACIÓN: Ratpack started for http://localhost:5050
```

We can now perform the following `POST` to get a fake token:

```bash
$> curl -X POST \
http://localhost:5050/ \
-H 'Content-Type: application/json' \
-d '{"jsonrpc":"2.0","id":1,"method":"UnicaOauth.1.getToken","params":{"params":[{"msisdn": "XXX","country": "XX"}],"gid":1009,"rid":2222}}'
```

**2. Mocking the Argentina API**: the script `scripts/argentina-api/argentina.groovy` launchs a server at `localhost:3001` that responds as the argentinian API does. Follow the instructions documented in the script in order to add new or modify existing JSON responses.

We aim to save in that folder curls and JSON responses from every integration with a new argentina endpoint; it'll act both as documentation and improvement of the mocked server.

```bash
$> cd [[Unica-Argentina]]/scripts/argentina-api-mock
$> groovy argentina.groovy
  ...
  [main] INFO ratpack.server.RatpackServer - Ratpack started (development) for http://localhost:3001
```

We can now perform a request to any of the endpoints that the script is configured; for example:

```bash
$> curl localhost:3001/telefonica/api/svagw/balancemanagement
```

**3. Local configuration**: the repository that contains the config files for local development is called `dev`. We need to tell our local server to route the Unica OAuth requests to `localhost:5050` and the third party requests to `localhost:3001`. For that purpose we have prepared a diff:

```bash
$> git clone https://git.code.tuenti.io/mvne/config/dev && cd dev # or git pull
$> git apply [[Unica-Argentina]]/scripts/local-dev/local-setup-mock-api.diff
```

**4. Checking everything is ok**: Once that both mock servers are up & runnign and our local config ready, we're ready to launch the `UnicaArgentina` server and perform requests.
You can launch the Tomcat server in the way you prefer (for example using an IDE) but here we provide a command line instruction to launch it from terminal.

```bash
$> cd [[Unica-Argentina]]
$> chmod +x scripts/local-dev/run.sh # needed only once
$> ./scripts/local-dev/run.sh
  ...
  INFORMACIÓN: Starting ProtocolHandler ["http-bio-8095"]
```

This is what we've got so far.
<p align="center">
  <img src="./doc/everything-keeps-local.png"/>
</p>

Now we can test both, the rest interface and the tService interface:

```bash
$> curl localhost:8095/UnicaArgentina/rest/customerAccounts/v2/accounts/1144347601/balance
```

```bash
$> curl http://localhost:8095/UnicaArgentina/ \
-H 'Content-Type: application/json' \
-d '{"jsonrpc":"2.0","id":1,"method":"MovistarArIntegrationService.1.getMultiplyBalanceOptionsForAccountType","params":{"params":["1128339358"],"gid":1009,"rid":2222}}'
```

#### Reaching Argentina API from local

**1. Running locally UnicaGateway**: run the Tomcat of a stable version of `UnicaGateway`. You can launch the server the way you prefer (for example using an IDE), here we've provided a maven command line script.

```bash
$> git clone https://git.code.tuenti.io/unica-gateway && cd unica-gateway # or git pull
$> chmod +x scripts/local-dev/run.sh # needed only the first time
$> ./scripts/local-dev/run.sh
  ...
  INFORMACIÓN: Starting ProtocolHandler ["http-bio-8080"]
```

**2. Reaching their API**: we need to create a tunnel through a dev machine. Run the following script:

```bash
$> cd [[Unica-Argentina]]
$> chmod +x scripts/local-dev/tunnel.sh # needed only the first time
$> ./scripts/local-dev/tunnel.sh
```

**3. Local configuration**: again, the repository that contains the config files for local development is `config/dev`. We need to tell our local server to route the Unica OAuth requests to `localhost:8080/UnicaGatewayService/`. For that purpose we have already prepared a diff:

```bash
$> git clone https://git.code.tuenti.io/mvne/config/dev && cd dev # or git pull
$> git apply [[Unica-Argentina]]/scripts/local-dev/local-setup-real-api.diff
```

**4. Checking everything is ok**: Once that we'll able to reach locally `UnicaGateway` and their production API, we're ready to launch the `UnicaArgentina` server and perform requests.
Again, you can launch the Tomcat server in the way you prefer (for example using an IDE) but if you wish in terminal:

```bash
$> cd [[Unica-Argentina]]
$> chmod +x scripts/local-dev/run.sh # needed only the first time
$> ./scripts/local-dev/run.sh
  ...
  INFORMACIÓN: Starting ProtocolHandler ["http-bio-8095"]
```

This is what we've got so far:
<p align="center">
  <img src="./doc/reaching-argentina-api-from-local.png"/>
</p>


Now we can test both, _unica_ requests (calling `UnicaGateway`) and the tService at `UnicaArgentina`

```bash
$> curl http://localhost:8080/UnicaGatewayService/ \
-H 'Content-Type: application/json' \
-d '{"jsonrpc":"2.0","id":1,"method":"UnicaCustomerAccounts.1.getAccountBalance","params":{"params":[{"country":"AR","msisdn":"1144347601"},"1144347601",{}],"gid":1009,"rid":2222}}'
```

```bash
$> curl http://localhost:8095/UnicaArgentina/ \
-H 'Content-Type: application/json' \
-d '{"jsonrpc":"2.0","id":1,"method":"MovistarArIntegrationService.1.getMultiplyBalanceOptionsForAccountType","params":{"params":["1128339358"],"gid":1009,"rid":2222}}'
```

### Integration and deployment

#### Git workflow

When your feature branch is ready, integrate it to `master` [flattening the git history](https://sites.google.com/a/tuenti.com/andres-viedma/java/flat-git-history).

**Notes about versioning**

We try to follow [Semantic Versioning](http://semver.org/):

* **Major** version: you probably won't code a new major version. Ask your team lead if in doubt.
* **Minor** version: you've declared a new function, dto class, value or constant in the IDL, needing to publish the IDL to flow.
* **Patch** version: Everything else.

:warning: [Flow](https://flow.tuenti.io) will automatically bump up the service version to the next **patch** when doing a release ticket (see **Deployment** section). However for **minor** version bumps you'll need to run the following command before integrating yout branch. This will change all `pom.xml` files for you.

```sh
mvn -B com.tuenti.maven.plugins:tuentiversions-maven-plugin:1.0.3:set-next-minor`
```

#### Deployment

##### Preproduction

If you just want to check your changes, you can deploy a snapshot of the service to a development environment.

The following command will create a [Jenkins](https://jenkins.tuenti.io) job that will generate a docker image tagged with the hash of the last commit (see `.ci.yml` file for details):

```
tu-ci snapshot
```

You'll be notified when the build is done. You can then deploy the image with `kubedeploy` or `kraken` (see [`k8s-definitions` documentation](https://doc.tuenti.io/repos/k8s-definitions/)).

##### Production

Once your code is properly integrated, you can deploy it to production.

Create a [Jira ticket](https://jira.tuenti.io/) of type **Automated Release** with the proper **Project** and **Component** fields (check list of [Jira automated projects in Flow](https://flow.tuenti.io/jira_comms/)). Also fill the **branch name** and **changeset** fields with your integrated commit.

Click on the **Start** button. The ticket will transition to **Prepared** status after a minute or so. Then click on the **Publish** button to transition it to **Published** status (it will also take some time). This will:

* Publish a non-snapshot version of your code in the images repository.
* Publish a tag of your version in the git repo.
* Merge two commits authored by Flow into `master`. These take care of bumping up the service version number to the next **patch** version. Remember that for **minor** versions you need to bump it up manually beforehand (see **Integration** section).

Once Flow is done, you can close the ticket. You'll now have to deploy the newly created image with `kubedeploy` or `kraken`, in a similar fashion to the preproduction instructions above (minus the `tu-ci` part).

### Verification

Movar team has some devices for testig purposes and some Argentina SIM numbers asigned.

The complete list of MSISDNs classified by OB is accessible in this document:
[SIMS test Tuenti-Novum](https://docs.google.com/spreadsheets/d/1QGr2-Ua1NvccnhqLv_XAVk1eyyIPwpBsBX7PuMS6WHE)
 - Check the pages "ARG B2B" & "ARG B2C"; you'll find the information divided for each team in Novum.
 - Search the **ACCOUNT/MOVAR** section in order to have the complete list of testing lines for each payment model (prepay, control, full).


## :white_check_mark: Developing tips

### Making a request to third party

If you need to make a request to the third party api, we suggest to use as template any of the `.curl` files located at `/scripts/argentina-api/`
You'll need a valid token & access to a miami dev machine.

In order to obtain an authentication token, one option is using the [fsm (production)](https://fsm.prd-mia.tuenti.io/)
 - Server: `Live`
 - Api: `UnicaOauth.1`
 - Config: `Movistar AR`
 - Method: `getToken`


## :books: Glossary

- **MSISDN(Mobile Station International Subscriber Directory Number):** Only applies for MOBILE. It is composed of 13 digits: The country prefix (54), a fixed 9, and the ANI. For example: 5491111111111
- **ANI:** 10 digit number that identifies the line. This applies for MOBILE, LANDLINE and INTERNET

---

> movar@tuenti.com
>
> diagrams has been created using Google Docs and originals are accesible for [edit](https://drive.google.com/drive/folders/1v5uHoXa_QJ-4VTbqMJoZUzjNGH7gi0Ii?usp=sharing)
