# Changelog

## [0.1.0](https://github.com/deepgram/spec-mock-go-sdk/compare/v0.0.1...v0.1.0) (2026-05-26)


### ⚠ BREAKING CHANGES

* **pkg:** wire Replace, drop internal stem fields and deprecated tier-2 fields ([#10](https://github.com/deepgram/spec-mock-go-sdk/issues/10))

### Features

* add Transport interface and built-in WebSocket subpackage ([53c17d2](https://github.com/deepgram/spec-mock-go-sdk/commit/53c17d2f38537ec12f217cb4e742326e6aec1bae))
* **api:** generated listen wire types from spec via smithy-go ShapeCodegenPlugin ([ba602d6](https://github.com/deepgram/spec-mock-go-sdk/commit/ba602d6140d39bca1d8f141fedc7a511c40f66fe))
* **api:** regen api/transport/http/rest.go with typed HTTPError ([ebb29d3](https://github.com/deepgram/spec-mock-go-sdk/commit/ebb29d3d8e341c18f0cf2d8361a13014f63cff80))
* **api:** regen api/types/listen_route.go with HTTPError.Typed decoder ([6f6743f](https://github.com/deepgram/spec-mock-go-sdk/commit/6f6743ff2b51519d85c1768f962f5128ea3ed266))
* generated transport packages + sagemakerruntime dep + 15 transport tests ([1d8a2a5](https://github.com/deepgram/spec-mock-go-sdk/commit/1d8a2a57a89c8219e9704410b447ab749b006e28))
* **listen:** idiomatic facade absorbs generated-type breaking changes ([3d2f3bc](https://github.com/deepgram/spec-mock-go-sdk/commit/3d2f3bc54e86ff16a570ccc898e3cb878f9146e7))
* **pkg/listen/rest:** route DoFile/DoStream/DoURL through httptransport.Invoke ([72f2c0b](https://github.com/deepgram/spec-mock-go-sdk/commit/72f2c0b12c16f24dfb6c8962f6ddbedc5c77f485))
* **pkg:** add requireFacadeOnly wire-helper for escape-hatch fields ([#15](https://github.com/deepgram/spec-mock-go-sdk/issues/15)) ([5040a94](https://github.com/deepgram/spec-mock-go-sdk/commit/5040a94e457fc8a487973ca8d32713764eb7e770))
* **pkg:** clean-slate facade reseed for listen WebSocket ([b5d7553](https://github.com/deepgram/spec-mock-go-sdk/commit/b5d75533e5de6bc962b254c1a3ba47fed072bc55))
* **pkg:** sagemaker routing options surface ([#18](https://github.com/deepgram/spec-mock-go-sdk/issues/18)) ([247dbf2](https://github.com/deepgram/spec-mock-go-sdk/commit/247dbf2a78d86f1e0b1cfa140c82a52ae2379556))
* **pkg:** tests, examples, and customer docs for new packages ([#20](https://github.com/deepgram/spec-mock-go-sdk/issues/20)) ([4f9fdff](https://github.com/deepgram/spec-mock-go-sdk/commit/4f9fdff98e9a407f162c8ba76f03b361163c482d))
* **pkg:** wire Replace, drop internal stem fields and deprecated tier-2 fields ([#10](https://github.com/deepgram/spec-mock-go-sdk/issues/10)) ([8535a6b](https://github.com/deepgram/spec-mock-go-sdk/commit/8535a6bab46be706113121c1cab18d55d0c36d97))
* **pkg:** wire SageMaker bidi runtime through live facade ([#19](https://github.com/deepgram/spec-mock-go-sdk/issues/19)) ([9e86f1d](https://github.com/deepgram/spec-mock-go-sdk/commit/9e86f1d4e49881ee4bf8b328295615048767598a))
* replace placeholder mock with deepgram-go-sdk, rewire listen REST ([784047c](https://github.com/deepgram/spec-mock-go-sdk/commit/784047c920a61a156fd14c7d240f1309ef2adb1c))
* scaffold mock fan-out target for deepgram/spec workflow ([3c6e2de](https://github.com/deepgram/spec-mock-go-sdk/commit/3c6e2de9ac61b0699ae694722ded197260eec27f))
* **ws:** production-resilience surface for the listen WebSocket facade ([6888264](https://github.com/deepgram/spec-mock-go-sdk/commit/6888264d3f3dd097f6835f4dde890050c3adfa81))


### Bug Fixes

* **examples:** synchronous prerecorded examples print transcript ([#23](https://github.com/deepgram/spec-mock-go-sdk/issues/23)) ([cb901ad](https://github.com/deepgram/spec-mock-go-sdk/commit/cb901ad9a5755af7d7eb777092ad818cf5017153))
* **pkg/listen/rest:** default empty Host to api.deepgram.com in baseURL ([a4408d1](https://github.com/deepgram/spec-mock-go-sdk/commit/a4408d128744b1e1b6c3822e38e913e6f30ad787))
* **pkg:** regen for audit-finding fixes ([#25](https://github.com/deepgram/spec-mock-go-sdk/issues/25)) ([4027355](https://github.com/deepgram/spec-mock-go-sdk/commit/4027355afb6bba2713b400ac4f384a9d2f498752))
* **pkg:** regen with SageMaker bindings repaired ([#21](https://github.com/deepgram/spec-mock-go-sdk/issues/21)) ([a58c346](https://github.com/deepgram/spec-mock-go-sdk/commit/a58c3467dfded2bea28581e680051558a8c6862f))
* **websocket:** regen api/ with gorilla/websocket; update ADR-0002 ([27c8176](https://github.com/deepgram/spec-mock-go-sdk/commit/27c817625920f2fd2f182d9efe78db3b7393c96c))

## Changelog

All notable changes are recorded here. Format follows [Keep a Changelog](https://keepachangelog.com/) and the project adheres to [Semantic Versioning](https://semver.org/).

This file is owned by [release-please](https://github.com/googleapis/release-please) — do not edit manually. Add changes via conventional commits.
