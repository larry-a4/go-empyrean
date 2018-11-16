// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{
	// Ethereum Foundation Go Bootnodes
	"enode://a979fb575495b8d6db44f750317d0f4622bf4c2aa3365d6af7c284339968eef29b69ad0dce72a4d8db5ebb4968de0e3bec910127f134779fbcb0cb6d3331163c@52.16.188.185:30303", // IE
	"enode://3f1d12044546b76342d59d4a05532c14b85aa669704bfe1f864fe079415aa2c02d743e03218e57a33fb94523adb54032871a6c51b2cc5514cb7c7e35b3ed0a99@13.93.211.84:30303",  // US-WEST
	"enode://78de8a0916848093c73790ead81d1928bec737d565119932b98c6b100d944b7a95e94f847f689fc723399d2e31129d182f7ef3863f2b4c820abbf3ab2722344d@191.235.84.50:30303", // BR
	"enode://158f8aab45f6d19c6cbf4a089c2670541a8da11978a2f90dbf6a502a4a3bab80d288afdbeb7ec0ef6d92de563767f3b1ea9e8e334ca711e9f8e2df5a0385e8e6@13.75.154.138:30303", // AU
	"enode://1118980bf48b0a3640bdba04e0fe78b1add18e1cd99bf22d53daac1fd9972ad650df52176e7c7d89d1114cfef2bc23a2959aa54998a46afcf7d91809f0855082@52.74.57.123:30303",  // SG

	// Ethereum Foundation C++ Bootnodes
	"enode://979b7fa28feeb35a4741660a16076f1943202cb72b6af70d327f053e248bab9ba81760f39d0701ef1d8f89cc1fbd2cacba0710a12cd5314d5e0c9021aa3637f9@5.1.83.226:30303", // DE
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
	"enode://30b7ab30a01c124a6cceca36863ece12c4f5fa68e3ba9b0b51407ccc002eeed3b3102d20a88f1c1d3c3154e2449317b8ef95090e77b312d5cc39354f86d5d606@52.176.7.10:30303",    // US-Azure geth
	"enode://865a63255b3bb68023b6bffd5095118fcc13e79dcf014fe4e47e065c350c7cc72af2e53eff895f11ba1bbb6a2b33271c1116ee870f266618eadfc2e78aa7349c@52.176.100.77:30303",  // US-Azure parity
	"enode://6332792c4a00e3e4ee0926ed89e0d27ef985424d97b6a45bf0f23e51f0dcb5e66b875777506458aea7af6f9e4ffb69f43f3778ee73c81ed9d34c51c4b16b0b0f@52.232.243.152:30303", // Parity
	"enode://94c15d1b9e2fe7ce56e458b9a3b672ef11894ddedd0c6f247e0f1d3487f52b66208fb4aeb8179fce6e3a749ea93ed147c37976d67af557508d199d9594c35f09@192.81.208.223:30303", // @gpip
}

// RinkebyBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network.
var RinkebyBootnodes = []string{
	"enode://a24ac7c5484ef4ed0c5eb2d36620ba4e4aa13b8c84684e1b4aab0cebea2ae45cb4d375b77eab56516d34bfbd3c1a833fc51296ff084b770b94fb9028c4d25ccf@52.169.42.101:30303", // IE
	"enode://343149e4feefa15d882d9fe4ac7d88f885bd05ebb735e547f12e12080a9fa07c8014ca6fd7f373123488102fe5e34111f8509cf0b7de3f5b44339c9f25e87cb8@52.3.158.184:30303",  // INFURA
	"enode://b6b28890b006743680c52e64e0d16db57f28124885595fa03a562be1d2bf0f3a1da297d56b13da25fb992888fd556d4c1a27b1f39d531bde7de1921c90061cc6@159.89.28.211:30303", // AKASHA
}

// GoerliBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Görli test network.
var GoerliBootnodes = []string{
	"enode://04fb7acb86f47b64298374b5ccb3c2959f1e5e9362158e50e0793c261518ffe83759d8295ca4a88091d4726d5f85e6276d53ae9ef4f35b8c4c0cc6b99c8c0537@40.70.214.166:40303",
	"enode://3d197d65ed92af6d0adf280ce486714fb641ef9f9f38f0bdd5ddd552666fc1132f033eb249a87f7f30086902c131f30f054f872ae80ac83eea6bd3760a7bbce2@40.70.214.166:30405",
	"enode://57f58f16fccdd9fb6f587565ac09af4b3b4b33d0fbd14252cc61d29a65b0d83c08419e67ac5292b9342090053526b847f2487278e609f4b4cd1dbf0f48105b2b@213.186.16.82:30303",
	"enode://5d9b1cba03738dfd23e12e4efb99b72623474fece2cc582c95e3ba7d481d519dea0029901f1f844116bab806044e8552f0431b21cf8d96010fc351b483330faa@13.78.10.94:30405",
	"enode://7592caf086d4d443905508492f40145bb1a0883ef7cbb9906b613eba6b501806e4ba0545a8e576236408e5b050e752e80a58445fb0ff2699b5ba4e334f481e40@13.78.10.94:30303",
	"enode://76850e0836d0074e060118bf57a627bfd8af3b59871fd16cb4d0ca826eda7a60b0e773f359335e5e3c6cea8a72b1efbf9a298a61b88d0c94ab1a6ea34f1d6c40@13.78.10.94:40303",
	"enode://87a7adc692793eb41918b74b7ba4aa9ec1b45a24917fd6e66118ffc9ffcac9d2672941fc10fd5a2d44e76d02628e273f861bf480311e31babd1ee211f5838e40@168.61.153.244:30405",
	"enode://9b1274fc252261bd9d8687bdc37cc3768551b93c9f3a3b3df2f4c7bbe6d797fb8c2ea6fb398114b2c6c6889a8257c244dfc57c1bb7b578c15cc5cc81fc0b3f79@168.61.153.244:30303",
	"enode://9bc25c32aaed85926a663563c8aa1c9abef6fb18e9282b7ae00584c9ed9ef8e353f18459b591c59b08f5f1ce692cf27cfdc5a0ff85312656aa65552e789f2315@40.74.91.252:40303",
	"enode://9ddf3e1ade168b2eea2d917dc32faffc727d53f488c78b293a523fda880bcca0b072506cf1ee6e743618d43f52e192fadc5ef5b43203a7f8e27b93a299248e3e@40.74.91.252:30405",
	"enode://a899e1b4551eb4d6e906a1313b8ba52e89eeb13412f1da058fd5a0cf261c235cb42fa38cc6c21b0fd5f5bcc5c5daa06945ea0410071cf34468a2f428454682ed@40.70.214.166:30303",
	"enode://d686ec8bf4bf0b205e8888e207352e6585232395e998cc1910b33be479c8405352ae1fc56aea79b2482b1b2d89412dc81091aa67ce775e335cd0f7d9dbfdfba3@84.196.20.71:30303",
	"enode://ea26ccaf0867771ba1fec32b3589c0169910cb4917017dba940efbef1d2515ce864f93a9abc846696ebad40c81de7c74d7b2b46794a71de8f95a0d019f494ff3@168.61.153.244:40303",
	"enode://ed70646a024612fa0db4fdb276a3add7ea322b13bec80dd1566186cc86cf9b853e1553eb0f49d3c4b9b37dc936f80e9ee0a2432b7b53f0b0d792fc2cdfb62861@88.19.163.180:30303",
	"enode://efaf6dad7a0773d911a6fcc44939faaba5d4802a7de8514bfacf9cc1ec9c292c82c1741eb4f14010895a273e7c94703cfc10c06068f6daa6ad25d4b0c0ca8e33@40.74.91.252:30303",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	"enode://06051a5573c81934c9554ef2898eb13b33a34b94cf36b202b69fde139ca17a85051979867720d4bdae4323d4943ddf9aeeb6643633aa656e0be843659795007a@35.177.226.168:30303",
	"enode://0cc5f5ffb5d9098c8b8c62325f3797f56509bff942704687b6530992ac706e2cb946b90a34f1f19548cd3c7baccbcaea354531e5983c7d1bc0dee16ce4b6440b@40.118.3.223:30304",
	"enode://1c7a64d76c0334b0418c004af2f67c50e36a3be60b5e4790bdac0439d21603469a85fad36f2473c9a80eb043ae60936df905fa28f1ff614c3e5dc34f15dcd2dc@40.118.3.223:30306",
	"enode://85c85d7143ae8bb96924f2b54f1b3e70d8c4d367af305325d30a61385a432f247d2c75c45c6b4a60335060d072d7f5b35dd1d4c45f76941f62a4f83b6e75daaf@40.118.3.223:30307",
}
