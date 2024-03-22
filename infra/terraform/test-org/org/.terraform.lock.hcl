# This file is maintained automatically by "terraform init".
# Manual edits may be lost in future updates.

provider "registry.terraform.io/hashicorp/external" {
  version     = "2.3.2"
  constraints = ">= 1.2.0, >= 2.2.2, < 3.0.0"
  hashes = [
    "h1:7F6FVQh7OcCgIH3YEJg1SJDSb1CU4qrCtGuI2EBHnL8=",
    "zh:020bf652739ecd841d696e6c1b85ce7dd803e9177136df8fb03aa08b87365389",
    "zh:0c7ea5a1cbf2e01a8627b8a84df69c93683f39fe947b288e958e72b9d12a827f",
    "zh:25a68604c7d6aa736d6e99225051279eaac3a7cf4cab33b00ff7eae7096166f6",
    "zh:34f46d82ca34604f6522de3b36eda19b7ad3be1e38947afc6ac31656eab58c8a",
    "zh:6959f8f2f3de93e61e0abb90dbec41e28a66daec1607c46f43976bd6da50bcfd",
    "zh:78d5eefdd9e494defcb3c68d282b8f96630502cac21d1ea161f53cfe9bb483b3",
    "zh:a81e5d65a343da9caa6f1d17ae0aced9faecb36b4f8554bd445dbd4f8be21ab6",
    "zh:b1d3f1557214d652c9120862ce27e9a7b61cb5aec5537a28240a5a37bf0b1413",
    "zh:b71588d006471ae2d4a7eca2c51d69fd7c5dec9b088315599b794e2ad0cc5e90",
    "zh:cfdaae4028b644dff3530c77b49d31f7e6f4c4e2a9e5c8ac6a88e383c80c9e9c",
    "zh:dbde15154c2eb38a5f54d0e7646bc67510004179696f3cc2bc1d877cecacf83b",
    "zh:fb681b363f83fb5f64dfa6afbf32d100d0facd2a766cf3493b8ddb0398e1b0f7",
  ]
}

provider "registry.terraform.io/hashicorp/google" {
  version     = "5.10.0"
  constraints = ">= 3.19.0, >= 3.39.0, >= 3.43.0, >= 3.45.0, >= 3.50.0, >= 3.53.0, >= 4.28.0, < 6.0.0"
  hashes = [
    "h1:3kD/GqYmZkA97ebToXu6qhhrwo+GQNmqq9xv30Qkhmw=",
    "zh:0f6a1feb5b3a128be6ef5fe0400ed800310a67e799c18aec7442161bb6d3ba36",
    "zh:13d591ba78e424c94ce5caaf176ab6b087b0e3af08a7b6bcd963673698cdefda",
    "zh:3bef54a2b24b06eef99f3df02e0fe4ac97f018c89f83e0faeb4ade921962565b",
    "zh:3f3755b8f5b9db1611d42a02c21f03c54577e4aad3cf93323792f131c671c050",
    "zh:61516eec734714ac48b565bee93cc2532160d1b4bd0320753799b829083b7060",
    "zh:9160848ad0b9becb522a0744dcb89474849906aa2436ed945c658fe201a724b0",
    "zh:aa5e79b01949cfedd874bf52958f90cf8f7d202600126c872127a9a156a3c17b",
    "zh:cef73a67031008b7d7ef3edfbcd5e1a9b04c0f2580d815401248025b741bc8e4",
    "zh:d2ad21ff9e9d2ad04146591c1b5784075e6df73e2bd243efd8d227d764b80b6e",
    "zh:f569b65999264a9416862bca5cd2a6177d94ccb0424f3a4ef424428912b9cb3c",
    "zh:f58b145081d20bce52e14bee0de73f5c018bc39b8c4736e23e1329df32f8bd45",
    "zh:fb82f6b5d1f992243ab8fe417659cdf9831202cf1e16fe7593d3967888b035cc",
  ]
}

provider "registry.terraform.io/hashicorp/google-beta" {
  version     = "5.10.0"
  constraints = ">= 3.19.0, >= 3.43.0, >= 3.50.0, >= 4.11.0, >= 4.28.0, < 6.0.0"
  hashes = [
    "h1:V+p4T5z6Mt6rFDupT/piBUFMg4zGJ153sawqNvqYWS0=",
    "zh:1004ac3733679254abcc7f5e9d594d9ee079cf071391a92f82b50077e07c70b5",
    "zh:1e25af33d20b6ab369860d5b7c746b4a3b3dccc061b14dde91b6ccccfe704cc4",
    "zh:2873a614a1dc1c460246edc95a558ad9befedf93490a0204bee8fb95362813cc",
    "zh:2f421e13247b3822ef3c2e07e1aee948116a5064c386466a53fb72486daded20",
    "zh:517c13cd146d3451789da8f13cbfa5355c3e88456cf762ad3918dada84a5f261",
    "zh:56553ae44f4089f5149551714daaf3c97205d4638dd93b0675ed777476d56048",
    "zh:6925a07bcb9ab70faa84bf36f87990025e3f9cd6c8cfab5260877f60086c8161",
    "zh:72454b65ee4a24896d215f7f7af41e31336865c86d6c20ea4acb63596e75ac0d",
    "zh:8b05f8a6ff51999bf65e3127618931647a00bc9abf739f0711151e4145cae3d5",
    "zh:a3b7d3b39740088174d121bc7e4e3ce27da0ebf0c87877f8fce9277b0046c75b",
    "zh:f569b65999264a9416862bca5cd2a6177d94ccb0424f3a4ef424428912b9cb3c",
    "zh:fe2af4fcda1b45d73ef8b8c728c150e00d1a4d5c0323b30d7d43c6f24ed78bcb",
  ]
}

provider "registry.terraform.io/hashicorp/kubernetes" {
  version     = "2.24.0"
  constraints = "~> 2.13"
  hashes = [
    "h1:u9lRMCdNXcB5/WQTZVMvGhNliW2pKOzj3SOVbu9yPpg=",
    "zh:0ed83ec390a7e75c4990ebce698f14234de2b6204ed9a01cd042bb7ea5f26564",
    "zh:195150e4fdab259c70088528006f4604557a051e037ebe8de64e92840f27e40a",
    "zh:1a334af55f7a74adf033eb871c9fe7e9e648b41ab84321114ef4ca0e7a34fba6",
    "zh:1ef68c3832691de21a61bf1a4e268123f3e08850712eda0b893cac908a0d1bc1",
    "zh:44a1c58e5a6646e62b0bad653319c245f3b635dd03554dea2707a38f553e4a52",
    "zh:54b5b374c4386f7f05b3fe986f9cb57bde4beab3bdf6ee33444f2b9a81b8af64",
    "zh:aa8c2687ab784b72f8cdad8d3c3673dea83b33561e7b3f2d287ef0d06ff2a9e5",
    "zh:e6ecba0503052ef3ad49ad56e17b2a73d9b55e30fcb82b040189d281e25e1a3b",
    "zh:f105393f6487d3eb1f1636ba42d10c82950ddfef852244c1bca8d526fa23a9a3",
    "zh:f17a8f1914ec66d80ccacecd40123362cf093abee3d3aa1ff9f8f687d8736f85",
    "zh:f394b12ef01fa0bdf666a43ad152eb3890134f35e635ea056b18771c292de46e",
    "zh:f569b65999264a9416862bca5cd2a6177d94ccb0424f3a4ef424428912b9cb3c",
  ]
}

provider "registry.terraform.io/hashicorp/null" {
  version     = "3.2.2"
  constraints = ">= 2.1.0, < 4.0.0"
  hashes = [
    "h1:zT1ZbegaAYHwQa+QwIFugArWikRJI9dqohj8xb0GY88=",
    "zh:3248aae6a2198f3ec8394218d05bd5e42be59f43a3a7c0b71c66ec0df08b69e7",
    "zh:32b1aaa1c3013d33c245493f4a65465eab9436b454d250102729321a44c8ab9a",
    "zh:38eff7e470acb48f66380a73a5c7cdd76cc9b9c9ba9a7249c7991488abe22fe3",
    "zh:4c2f1faee67af104f5f9e711c4574ff4d298afaa8a420680b0cb55d7bbc65606",
    "zh:544b33b757c0b954dbb87db83a5ad921edd61f02f1dc86c6186a5ea86465b546",
    "zh:696cf785090e1e8cf1587499516b0494f47413b43cb99877ad97f5d0de3dc539",
    "zh:6e301f34757b5d265ae44467d95306d61bef5e41930be1365f5a8dcf80f59452",
    "zh:78d5eefdd9e494defcb3c68d282b8f96630502cac21d1ea161f53cfe9bb483b3",
    "zh:913a929070c819e59e94bb37a2a253c228f83921136ff4a7aa1a178c7cce5422",
    "zh:aa9015926cd152425dbf86d1abdbc74bfe0e1ba3d26b3db35051d7b9ca9f72ae",
    "zh:bb04798b016e1e1d49bcc76d62c53b56c88c63d6f2dfe38821afef17c416a0e1",
    "zh:c23084e1b23577de22603cff752e59128d83cfecc2e6819edadd8cf7a10af11e",
  ]
}

provider "registry.terraform.io/hashicorp/random" {
  version     = "3.6.0"
  constraints = ">= 2.1.0, >= 2.2.0, >= 2.3.1, < 4.0.0"
  hashes = [
    "h1:R5Ucn26riKIEijcsiOMBR3uOAjuOMfI1x7XvH4P6B1w=",
    "zh:03360ed3ecd31e8c5dac9c95fe0858be50f3e9a0d0c654b5e504109c2159287d",
    "zh:1c67ac51254ba2a2bb53a25e8ae7e4d076103483f55f39b426ec55e47d1fe211",
    "zh:24a17bba7f6d679538ff51b3a2f378cedadede97af8a1db7dad4fd8d6d50f829",
    "zh:30ffb297ffd1633175d6545d37c2217e2cef9545a6e03946e514c59c0859b77d",
    "zh:454ce4b3dbc73e6775f2f6605d45cee6e16c3872a2e66a2c97993d6e5cbd7055",
    "zh:78d5eefdd9e494defcb3c68d282b8f96630502cac21d1ea161f53cfe9bb483b3",
    "zh:91df0a9fab329aff2ff4cf26797592eb7a3a90b4a0c04d64ce186654e0cc6e17",
    "zh:aa57384b85622a9f7bfb5d4512ca88e61f22a9cea9f30febaa4c98c68ff0dc21",
    "zh:c4a3e329ba786ffb6f2b694e1fd41d413a7010f3a53c20b432325a94fa71e839",
    "zh:e2699bc9116447f96c53d55f2a00570f982e6f9935038c3810603572693712d0",
    "zh:e747c0fd5d7684e5bfad8aa0ca441903f15ae7a98a737ff6aca24ba223207e2c",
    "zh:f1ca75f417ce490368f047b63ec09fd003711ae48487fba90b4aba2ccf71920e",
  ]
}

provider "registry.terraform.io/hashicorp/time" {
  version     = "0.10.0"
  constraints = ">= 0.5.0"
  hashes = [
    "h1:EeF/Lb4db1Kl1HEHzT1StTC7RRqHn/eB7aDR3C3yjVg=",
    "zh:0ab31efe760cc86c9eef9e8eb070ae9e15c52c617243bbd9041632d44ea70781",
    "zh:0ee4e906e28f23c598632eeac297ab098d6d6a90629d15516814ab90ad42aec8",
    "zh:3bbb3e9da728b82428c6f18533b5b7c014e8ff1b8d9b2587107c966b985e5bcc",
    "zh:6771c72db4e4486f2c2603c81dfddd9e28b6554d1ded2996b4cb37f887b467de",
    "zh:78d5eefdd9e494defcb3c68d282b8f96630502cac21d1ea161f53cfe9bb483b3",
    "zh:833c636d86c2c8f23296a7da5d492bdfd7260e22899fc8af8cc3937eb41a7391",
    "zh:c545f1497ae0978ffc979645e594b57ff06c30b4144486f4f362d686366e2e42",
    "zh:def83c6a85db611b8f1d996d32869f59397c23b8b78e39a978c8a2296b0588b2",
    "zh:df9579b72cc8e5fac6efee20c7d0a8b72d3d859b50828b1c473d620ab939e2c7",
    "zh:e281a8ecbb33c185e2d0976dc526c93b7359e3ffdc8130df7422863f4952c00e",
    "zh:ecb1af3ae67ac7933b5630606672c94ec1f54b119bf77d3091f16d55ab634461",
    "zh:f8109f13e07a741e1e8a52134f84583f97a819e33600be44623a21f6424d6593",
  ]
}

provider "registry.terraform.io/integrations/github" {
  version     = "6.0.0"
  constraints = "~> 6.0"
  hashes = [
    "h1:9pdD0wlgzxXpJt41zvAPTBJlSVQjrivGj/PKXonvjdI=",
    "h1:C0tQYTi4xfFbv49ohtqcnUS6N3zkieMUlGgVtZa2KNg=",
    "h1:IsJlhZqzDak5PE4u/DGPpVuh007NWn6RGL42sZNUZtE=",
    "h1:JN9FDT93mtIFE9oTZFJN8iBBkYM4VUBN35H9ejI0pKA=",
    "h1:KXepSQ13ED8xN5b74H4KPbWkm03U53F5ey+Htk+SLlk=",
    "h1:LeDpsKXQvLh5IPHNj0i5/2j0G3QocWKosoI4Vt+R2GY=",
    "h1:YGUmIpK8zBYXssW5vJcQEaRvimE/kxcIifcEDnfThMQ=",
    "h1:aPVKHd7sUHWfHO2nY2jUaEGa6YK89KPNCykbWKVMPoA=",
    "h1:cHq9ip3mg1zIEXQPvi7Zqb2dQcsQBZlthQWFMaRzbUE=",
    "h1:gXBHc4e5JRRkX35POfXbiuBBPNxjQ6KAH0d9QR+jeWk=",
    "h1:hgxwIjasPR+EYjFCms6PkTBYYx1+VJJfcOQ5/UjTReQ=",
    "h1:jLOsi4Qu5g5D2/n/xg/CljAKCRH9F9paiWRZtyzWR+k=",
    "h1:mTMlNk78lzcXTm4kgqCVFETkGIt0RgwQGzwDvI5HjbQ=",
    "h1:p86Nsa+Fo29MVOXYG8x+M0MP/R4QS8nl4CUu6Us+/YU=",
    "zh:0d12fde69c54d358af3a45cf1610b711e1cd6a5d0be8d71c24729f28faa4a67d",
    "zh:501fd9a181bbb1f3e70c3a54463bc16974569dddd1311fdd682c3b893ebc8455",
    "zh:69a486e2b2db2f7ff947027e5e245b48a1f71e10955e7243419c15d9d8330d54",
    "zh:6f45927d00337db1ebce1da51be1033ffb632b470f901698e12cd51e1d2e16db",
    "zh:7734fccb5594f72d8f0bd501f83bbdcc8cf69df27b54631a24d27cff5cead9ab",
    "zh:77ffc5ec11754a1c94af468c17a95409a36b2696a4d1f656cef893d931d20b2e",
    "zh:79dc9da6aff69825e66869ccdff83ce453f1374ee08152645bc324885e1d1b42",
    "zh:7cbf2e8a01133b4ad442854b4baceb97cf4b9f43d69684a35583eaeb998cbc5d",
    "zh:84d9788a46f57572a348a52bbbcb347786d967d2169825587b7bf1fc6d052d71",
    "zh:a0b89fcce44c397c5f61286351f1752c154e1f238555c4a69a6cd49a57f79d02",
    "zh:c2fe95b549239b01ae7956f00279ab6653521843b7009231aec3eb898c8dc395",
    "zh:c58aa97ae9b24c260f1f8c7a4c2a7ffc75fe0c2ffa0cb9986d99e855c11a0cbb",
    "zh:cb32a38fb412935ee021f70050db07dbe9ec698bcf149275fd3381565eb9b5d1",
    "zh:d6ae9b8fa87f3fefe13976504731f661fd93729f8167dd5b7056b1d325b745e4",
  ]
}
