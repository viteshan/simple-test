package stellar

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/viteshan/go/build"

	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

func TestGenerateKey(t *testing.T) {
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(pair.Seed())
	// SA7LTOMFIC5TT4RY5VQSMCHGA642OEQZOH4PMWXW6NDGNH7JVFDSER7I
	log.Println(pair.Address())
	// GC7E52HMEVSLXRR4USWGFZA7F3F5I2MFUUWDKVAU6LAMZOHOZIDFNAR6
}

func TestCreateAccount(t *testing.T) {
	// pair is the pair that was generated from previous example, or create a pair based on
	// existing keys.
	kp, e := keypair.Parse("GC7E52HMEVSLXRR4USWGFZA7F3F5I2MFUUWDKVAU6LAMZOHOZIDFNAR6")
	if e != nil {
		panic(e)
	}
	address := kp.Address()
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + address)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
	/**
	{
		"_links": {
		"transaction": {
			"href": "https://horizon-testnet.stellar.org/transactions/96603e7fd5fae4d6c7a9e1ca71a00b34ddd2200a4e3ff1c924ecc83d3e82eac7"
		}
	},
		"hash": "96603e7fd5fae4d6c7a9e1ca71a00b34ddd2200a4e3ff1c924ecc83d3e82eac7",
		"ledger": 998956,
		"envelope_xdr": "AAAAABB90WssODNIgi6BHveqzxTRmIpvAFRyVNM+Hm2GVuCcAAAAZAAAALsAAligAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAvk7o7CVku8Y8pKxi5B8uy9RphaUsNVQU8sDMuO7KBlYAAAAXSHboAAAAAAAAAAABhlbgnAAAAEBKaQ/WPl8RGkHHiSNaDuBRrLArs1M3SitYKwgPQxdXLZqzfIyTkWDguORmBwa2+DlOVkJaGxBOR28DdyzluLkB",
		"result_xdr": "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=",
		"result_meta_xdr": "AAAAAQAAAAIAAAADAA8+LAAAAAAAAAAAEH3Rayw4M0iCLoEe96rPFNGYim8AVHJU0z4ebYZW4JwAHFsaEcZwtAAAALsAAlifAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAA8+LAAAAAAAAAAAEH3Rayw4M0iCLoEe96rPFNGYim8AVHJU0z4ebYZW4JwAHFsaEcZwtAAAALsAAligAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAAADz4sAAAAAAAAAAC+TujsJWS7xjykrGLkHy7L1GmFpSw1VBTywMy47soGVgAAABdIdugAAA8+LAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMADz4sAAAAAAAAAAAQfdFrLDgzSIIugR73qs8U0ZiKbwBUclTTPh5thlbgnAAcWxoRxnC0AAAAuwACWKAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEADz4sAAAAAAAAAAAQfdFrLDgzSIIugR73qs8U0ZiKbwBUclTTPh5thlbgnAAcWwLJT4i0AAAAuwACWKAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA=="
	}
	*/
}

func TestGetBalance(t *testing.T) {
	kp, e := keypair.Parse("GC7E52HMEVSLXRR4USWGFZA7F3F5I2MFUUWDKVAU6LAMZOHOZIDFNAR6")
	if e != nil {
		panic(e)
	}
	account, err := horizon.DefaultTestNetClient.LoadAccount(kp.Address())
	if err != nil {
		log.Fatal(err)
	}

	build.Payment()
	fmt.Println("Balances for account:", kp.Address())

	for _, balance := range account.Balances {
		log.Println(balance)
	}
}
