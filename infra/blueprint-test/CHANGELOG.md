# Changelog

## [0.7.0](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/compare/infra/blueprint-test/v0.6.1...infra/blueprint-test/v0.7.0) (2023-07-20)


### Features

* Add HTTP Assert test helpers ([#1707](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1707)) ([9c423f9](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/9c423f910c14899eb311bf9b026439eb70378602))
* add retry for kpt commands ([#1717](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1717)) ([55c9c8d](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/55c9c8dcf85b8eacdb8ff2c9d19582a445e192ab))

## [0.6.1](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/compare/infra/blueprint-test/v0.6.0...infra/blueprint-test/v0.6.1) (2023-06-27)


### Bug Fixes

* blueprint-test tests ([#1675](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1675)) ([6ed5385](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/6ed538548fb9fd91a81663796efecb5e53c8a66e))
* update bpt go modules ([#1677](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1677)) ([4c9aaec](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/4c9aaeca68db7198165d52227b3d03752d8f817d))

## [0.6.0](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/compare/infra/blueprint-test/v0.5.2...infra/blueprint-test/v0.6.0) (2023-06-13)


### Features

* update to bpt GO 1.20 and rework krm test ([#1619](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1619)) ([50c2ab3](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/50c2ab3165ab5eb159a8569ec90cd1518d427b7c))

## [0.5.2](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/compare/infra/blueprint-test/v0.5.1...infra/blueprint-test/v0.5.2) (2023-05-11)


### Bug Fixes

* bump GO modules and address lint ([#1541](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1541)) ([6b76dc1](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/6b76dc17db4e64a6aff52b980d5c3ac01b2a901a))

## [0.5.1](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/compare/infra/blueprint-test/v0.5.0...infra/blueprint-test/v0.5.1) (2023-03-20)


### Bug Fixes

* kpt tests without existing working dir ([#1447](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1447)) ([c9cc7af](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/c9cc7af901d8ff6c100358c540eb9eea5f8015a4))

## [0.5.0](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/compare/infra/blueprint-test/v0.4.1...infra/blueprint-test/v0.5.0) (2023-02-28)


### Features

* update blueprint-test to GO 1.18 and test fixes ([#1373](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1373)) ([0234ad6](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/0234ad6f0da169aec58a9fd848094907aa4b6851))


### Bug Fixes

* **deps:** update go modules ([#1416](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1416)) ([5f01e1f](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/5f01e1ffd04d9ad47bf7bdb28c92716028d1977f))
* update blueprint-test for kpt v1.0.0-beta.16+ ([#1367](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1367)) ([3613491](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/36134916e2fd859b0aea4384c1b4a5ab79d65eac))

## [0.4.1](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/compare/infra/blueprint-test/v0.4.0...infra/blueprint-test/v0.4.1) (2023-01-10)


### Bug Fixes

* **deps:** update for go-sdk refactor ([#1217](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1217)) ([5c50728](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/5c50728b825fda6187ca9b73151741c733e623ec))
* remove terraform plan file needed for the terraform vet execution after validation ([#1321](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1321)) ([1edc5df](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/1edc5df7267c78a917eec0f2b5ad3f4024ca5e98))

## [0.4.0](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/compare/infra/blueprint-test/v0.3.0...infra/blueprint-test/v0.4.0) (2022-12-02)


### Features

* allow var overrides for workspace mode ([#1292](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1292)) ([f6ffa1f](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/f6ffa1f60039f03a6fb77e122894641caa739fef))
* enable no color ([#1293](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1293)) ([06fae23](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/06fae232e1f97b1d78df6809eff65898fddb5268))
* new test strategy for redeploy validation ([#1286](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1286)) ([de5d509](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/de5d5090980f5f12e0321365a935e119493518ec))

## [0.3.0](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/compare/infra/blueprint-test/v0.2.0...infra/blueprint-test/v0.3.0) (2022-08-30)


### Features

* add project ID param to tft vet ([#1226](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1226)) ([e95dc64](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/e95dc64d9b4596135cfc8bac481402c739e1c6a4))
* blueprint-test file logger ([a284f16](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/a284f164fccc58db86dfb8999b8013642e5d2bd7))

## [0.2.0](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/compare/infra/blueprint-test/v0.1.0...infra/blueprint-test/v0.2.0) (2022-08-03)


### Features

* add support for retryable tf errors ([#1198](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1198)) ([bcf67d6](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/bcf67d6d5aa193077c961c529f14df56e80b9e7a))
* Add support for terraform vet in blueprint test ([#1191](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1191)) ([e3179df](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/e3179dfc63abf2fd7cf291b24531abbe3cba02ff))
* expose setup outputs via tft ([#1203](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1203)) ([4ea786f](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/4ea786f947dfa58799f5e4736e511ca59668958b))

## [0.1.0](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/compare/infra/blueprint-test/v0.0.1...infra/blueprint-test/v0.1.0) (2022-06-13)


### Features

* add support for backend configuration tf blueprint test ([#1165](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1165)) ([442b49e](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/442b49ebe347d2415840967200d280bdf590cbe1))

## [0.0.1](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/compare/infra/blueprint-test-v0.0.1...infra/blueprint-test/v0.0.1) (2022-06-07)


### Features

* add getter for krmt build directory ([#1106](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1106)) ([fd68a6b](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/fd68a6bdc9a90d0f340fdad80bfcfc8119137a0f))
* add golden string sanitizer option ([#1109](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1109)) ([0e962c6](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/0e962c6ff0f5fa4f38cab62c31e78bd67d5923de))
* add goldenfile sanitizers ([#1074](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1074)) ([c98be35](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/c98be3561409b0051a1a5b2502eb603766d2c4a5))
* add helper for goldenfiles ([#1067](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1067)) ([1bf5397](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/1bf53970d457786fa2b3dc79a42b77887d1c7bb5))
* add KRM benchmark example ([#982](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/982)) ([6854aa0](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/6854aa08ed6f5edeba8884aec1d89745d1f64a2b))
* add support for runf in gcloud blueprint test ([#1070](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1070)) ([3842083](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/3842083683a3218d7864efa5e545dc4958cc3ecb))
* add test result utils ([#1005](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1005)) ([608c349](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/608c349bf8e4b68bb1f211094de5e8c91f881521))
* add test yaml config ([#986](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/986)) ([fe03487](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/fe034876e2e780bce0906252026115c851580f7a))
* add the first draft of the user guide for the test framework ([#983](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/983)) ([5dcd154](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/5dcd1546f9d5ec5ab39743e5181feffb8877a1ea))
* add transform func for removing a json value by path ([#1110](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1110)) ([9f9a444](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/9f9a444009fd9a35d074efff5803d2a2fb8572e8))
* adds logic for copying additional resources to the build directory ([#1118](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1118)) ([8383c92](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/8383c92c8a54322eebca8a058550b46396f043aa))
* **cli:** Allow customization of setup sa key ([#1065](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1065)) ([7c9f83c](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/7c9f83caf2fe77b69dd4af8bef6c3496b14d3af2))
* export GetTFOptions method for tft ([#1003](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1003)) ([5e783cf](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/5e783cf7c716104ad113064d9e1b9aea4dc7a999))
* initialize KRM blueprint testing ([#977](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/977)) ([2953e46](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/2953e46e28f4085c733243e3a4914f52aa105f2e))
* initialize terraform blueprint testing ([#945](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/945)) ([723b19c](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/723b19ce02d0e04e1f12f117aa2fe9ba44cad5e5))
* remove list-setter dep for kpt ([#1088](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1088)) ([bad09af](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/bad09af2b45598ca08990d3bcb722560b661b3e0))
* support complex type setup outputs ([#997](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/997)) ([39b4ef0](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/39b4ef08bb23b352ea6ff7073f942ec0b5a50fc7))
* **tft:** activate creds from setup ([#1062](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1062)) ([08c972c](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/08c972c3768cae717df3f33a43785bf21b183a13))


### Bug Fixes

* **bptest:** compute config discovery from int test dir ([#1025](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1025)) ([bea525b](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/bea525b1cf5203f522bd0e9f42bce45605885c41))
* bumb the gjson version to the latestg ([#1011](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1011)) ([2c665e7](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/2c665e7bd84a189f225fbbf42f4ea5d0b69fa42a))
* **krm-test:** add option for custom setters ([#981](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/981)) ([78afb5d](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/78afb5d4cfd83922d82980a6493bd0a7dab78e12))
* mark tests as skipped due to test config ([#1063](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1063)) ([0687139](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/068713996f3641114bf1fed1937d4cec09ddc3f5))
* recognize prow PR commit ([#993](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/993)) ([e8c47de](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/e8c47de6a66b1dde620da57ecd752f59de32b7f4))
* upgrades kpt-sdk dependency to remove the gakekeeper lib reference ([#1090](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/1090)) ([727d5c1](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/commit/727d5c1b1fafbd45ddfaea7cad99da379891fc6e))
