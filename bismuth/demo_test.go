package bismuth

import (
	"fmt"
	"testing"

	"github.com/viteshan/simple-test/bismuth/native"
)

func TestBismuth(t *testing.T) {
	const nodeServer string = "127.0.0.1"
	const nodePort int = 25658
	//const nodePort int = 18115

	connectionBismuthAPI := native.New(nodeServer, nodePort)
	defer connectionBismuthAPI.Close()

	//// Ask for general status
	//println("statusjson:")
	//ret = connectionBismuthAPI.Command("statusjson")
	//println(ret)
	//println("")

	//commandAndPrintln(connectionBismuthAPI.Command, "balancegetjson", "f6c0363ca1c5aa28cc584252e65a63998493ff0a5ec1bb16beda9bac")
	//commandAndPrintln(connectionBismuthAPI.Command, "balanceget", "f6c0363ca1c5aa28cc584252e65a63998493ff0a5ec1bb16beda9bac")
	//commandAndPrintln(connectionBismuthAPI.Command, "blockgetjson", "1270871")
	//commandAndPrintln(connectionBismuthAPI.Command, "blockgetjson", "1270870")
	commandAndPrintln(connectionBismuthAPI.Command, "blockgetjson", "1275725")
	commandAndPrintln(connectionBismuthAPI.Command, "addlistlimjson", "95bf1ed5eba5dcaaa4a2c4e95617b6f48106ff692377a4fc149391b8", 10)
	//commandAndPrintln(connectionBismuthAPI.Command, "api_getblocksince", 1271360)

	// Gets a specific block details

	//ret = connectionBismuthAPI.Command("blockget", 558742)
	//println(ret)
	//println("blockget(558742):", ret)
	//println("")
	//
	//// Gets latest block and diff
	//println("difflast:")
	//ret = connectionBismuthAPI.Command("difflast")
	//println(ret)
	//println("")
	//
	//// Gets tx detail, raw format
	//println("api_gettransaction('K1iuKwkOac4HSuzEBDxmqb5dOmfXEK98BaWQFHltdrbDd0C5iIEbh/Fj', false):")
	//ret = connectionBismuthAPI.Command("api_gettransaction", "K1iuKwkOac4HSuzEBDxmqb5dOmfXEK98BaWQFHltdrbDd0C5iIEbh/Fj", false)
	//println(ret)
	//println("")
	//
	//// Gets tx detail, json format
	//println("api_gettransaction('K1iuKwkOac4HSuzEBDxmqb5dOmfXEK98BaWQFHltdrbDd0C5iIEbh/Fj', true):")
	//ret = connectionBismuthAPI.Command("api_gettransaction", "K1iuKwkOac4HSuzEBDxmqb5dOmfXEK98BaWQFHltdrbDd0C5iIEbh/Fj", true)
	//println(ret)
	//println("")
	//
	//// Gets addresses balances
	//println("api_listbalance(['731337bb0f76463d578626a48367dfea4c6efcfa317604814f875d10','340c195f768be515488a6efedb958e135150b2ef3e53573a7017ac7d'], 0, true):")
	//addresses := []string{"731337bb0f76463d578626a48367dfea4c6efcfa317604814f875d10", "340c195f768be515488a6efedb958e135150b2ef3e53573a7017ac7d"}
	//ret = connectionBismuthAPI.Command("api_listbalance", addresses, 0, true)
	//println(ret)
	//println("")
	//
	//// Ask for a new keys/address set
	//println("keygen:")
	//ret = connectionBismuthAPI.Command("keygen")
	//println(ret)
	//println("")
}

func commandAndPrintln(fn func(...interface{}) string, cmds ...interface{}) {
	fmt.Printf("%s\t   %s\n", cmds, fn(cmds...))
}

func TestShell(t *testing.T) {
	sbp := []string{"vite_0307e7f8f7dd7a5566c8fd5ba583b6354bc2d5bf0c91f5957c",
		"vite_047c39256cb6d80eaa7aef39c8c60a39c38da69ea776a1af58",
		"vite_0ab98c623471d3512bb6cbe06c243e977ab6e91c6abffd914e",
		"vite_10513d54e0c38a304ad9e7902c82277328b4df76dd31871f37",
		"vite_1164f7ddda038466e510e26a6feba9214f6e7e2643baa4f016",
		"vite_12bd06e0318062feaa343dcd48ed783aa902e2abbfeed6a26e",
		"vite_165a295e214421ef1276e79990533953e901291d29b2d4851f",
		"vite_1d46ef813f218a2c7fab8b4a71b6176220a9c8a7a417dee341",
		"vite_1d656d6796c6894d18c1a64ac2b78e6f082e76808572578e8f",
		"vite_1dfe68a6e6f3d7970cad2d0b92d21c23fdf6c847b28d7f2f89",
		"vite_33916a1b849611e0d97d6b335d19527b01c1421c305d4054f7",
		"vite_35160e6804adafffb51dd5569c35d6f9f25faea43292f3b0e7",
		"vite_3ae6649db0f94239f5256c4bcd3b10c24738fbe5d24c00f809",
		"vite_420baeefd649352532e79fabac194f4f95725ab976396fc36c",
		"vite_4a04e8b6ba89f741c6945868a010ba82b22ca993ab89c6225d",
		"vite_4bcee903088951aeb55f16bfa5aade79cbe93d7d98216b3e86",
		"vite_5809101f1839992b53802d9c0eb9c0373febcc18289c736c19",
		"vite_62c8cbd276de2d086d48ba15b8b20ad8096021daa7a7f622cb",
		"vite_6e2e54c8c3f620a31a358650b3d92973c9d729b300a427ba52",
		"vite_753f84104dae9858a51173855db21463f23a00753a3491885b",
		"vite_785c872f703666b6088767d9809fa8aaf932d0e385f02a0c92",
		"vite_7adabd80763b42543078bcfd24888489fb790f159dc1129aef",
		"vite_7ba4daf7b6902eb4a1c00bb97500256cf82b97e4cf827f5ed4",
		"vite_7c655aca252d6f8b3788449bb39b330b32397b549b0a2d1000",
		"vite_82d129931823a8a47275baf4c5c853713716800220b2ceaafd",
		"vite_8370865362e739fb71615b8b33f9e394d85743093bdfaede6c",
		"vite_9065ff0e14ebf983e090cde47d59fe77d7164b576a6d2d0eda",
		"vite_94badf80abab06dc1cdb4d21038a6799040bb2feb154f730cb",
		"vite_995769283a01ba8d00258dbb5371c915df59c8657335bfb1b2",
		"vite_9b33ce9fc70f14407db75cfa8453680f364e6674c7cc1fb785",
		"vite_9d8b73d3408ccbf1f6efce89a1563660c5b1c16aa3dec131e3",
		"vite_a1ffbed38cfcdb0f4bb0b79168c1b7970e774120d351aed4be",
		"vite_aa7f76480db9e4072231d52f7b7abcccc8197b217a2dd5e818",
		"vite_aa8c1fa66636479ab725e2f32f9b0e17df76ae0b6adfc7675f",
		"vite_aadf6275ecbf07181c5c37c7d709aebf06553e470345d3f699",
		"vite_adee0a2e60c50273ef54a042593628b590564ca921d95e73b8",
		"vite_af4e623adc61ab661fd2c5dbbce23b12eaadc92153c032cd64",
		"vite_b6fe7fb14765d8cc15b72816dcd1bc058eab07765edead449b",
		"vite_b9351a94769d85453d7ebd8dce868e5eaa56c7c5a23216ce92",
		"vite_babf3d7d9e454a8af9d80ab828407600546798c42370704efb",
		"vite_c1d11e6eda9a9b80e388a38e0ac541cbc3333736233b4eaaab",
		"vite_c38760d672f05ce5e1aef97884ecf3085b9ebe10000be47fb7",
		"vite_cdd075a61de71a076c94bf4fd45e5b6bc3c48ac65db912bc57",
		"vite_d12cfc15515e56d289ff0d17dadc10e1ceeca9129063119a80",
		"vite_d4672160eabc76770cb0cf8c4fbc8bc3004c2a5c84829007b8",
		"vite_da12e0c03bf219395c13dc0367ac3944b1de09b30033a13c36",
		"vite_dff5ee13c87ed2f205ef87d820b3cd8e97c181b1bb6781c602",
		"vite_e79efe2964388d40b75dc44f6b3e05956c5f4c90f8a6263f98",
		"vite_e7a01e66d920c6c5ce82e9353ca57267f6534e357bbee63063",
		"vite_ff3a9e77549c2f084bc8ce0f97d964ba607d36f9594e9cdfb1",
	}

	v2_3_1 := []string{
		"vite_9d8b73d3408ccbf1f6efce89a1563660c5b1c16aa3dec131e3",
		"vite_a1ffbed38cfcdb0f4bb0b79168c1b7970e774120d351aed4be",
		"vite_d12cfc15515e56d289ff0d17dadc10e1ceeca9129063119a80",
		"vite_e7a01e66d920c6c5ce82e9353ca57267f6534e357bbee63063",
		"vite_165a295e214421ef1276e79990533953e901291d29b2d4851f",
		"vite_10513d54e0c38a304ad9e7902c82277328b4df76dd31871f37",
		"vite_da12e0c03bf219395c13dc0367ac3944b1de09b30033a13c36",
		"vite_6e2e54c8c3f620a31a358650b3d92973c9d729b300a427ba52",
		"vite_82d129931823a8a47275baf4c5c853713716800220b2ceaafd",
		"vite_1d46ef813f218a2c7fab8b4a71b6176220a9c8a7a417dee341",
		"vite_b9351a94769d85453d7ebd8dce868e5eaa56c7c5a23216ce92",
		"vite_62c8cbd276de2d086d48ba15b8b20ad8096021daa7a7f622cb",
		"vite_cdd075a61de71a076c94bf4fd45e5b6bc3c48ac65db912bc57",
		"vite_33916a1b849611e0d97d6b335d19527b01c1421c305d4054f7",
		"vite_1dfe68a6e6f3d7970cad2d0b92d21c23fdf6c847b28d7f2f89",
		"vite_12bd06e0318062feaa343dcd48ed783aa902e2abbfeed6a26e",
		"vite_dff5ee13c87ed2f205ef87d820b3cd8e97c181b1bb6781c602",
		"vite_4a04e8b6ba89f741c6945868a010ba82b22ca993ab89c6225d",
		"vite_35160e6804adafffb51dd5569c35d6f9f25faea43292f3b0e7",
		"vite_785c872f703666b6088767d9809fa8aaf932d0e385f02a0c92",
		"vite_7adabd80763b42543078bcfd24888489fb790f159dc1129aef",
		"vite_1164f7ddda038466e510e26a6feba9214f6e7e2643baa4f016",
		"vite_7c655aca252d6f8b3788449bb39b330b32397b549b0a2d1000",
	}

LL:
	for _, v := range sbp {
		for _, vv := range v2_3_1 {
			if vv == v {
				continue LL
			}
		}
		println(v)
	}

}
