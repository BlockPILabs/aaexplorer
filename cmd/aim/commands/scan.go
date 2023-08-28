package commands

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	aimos "github.com/BlockPILabs/aa-scan/internal/os"
	"github.com/BlockPILabs/aa-scan/task"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/procyon-projects/chrono"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// ScanCmd ...
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan block",
	Run: func(cmd *cobra.Command, args []string) {
		taskScheduler := chrono.NewDefaultTaskScheduler()
		// db start
		err := entity.Start(config)
		if err != nil {
			return
		}
		_, err = taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
			//parser.ScanBlock()
		}, 5*time.Second)

		if err == nil {
			log.Print("Task: scan block has been scheduled successfully.")
		}
		//test1()
		//testSubscribe()
		//getAbi()
		//moralis.TestToken("0x3a55815977ab0e12E4Fcf1a66165142C41dbda26")
		test2()
		task.AssetSync()
		aimos.TrapSignal(logger, func() {})

		//task.InitTask()

		// Run forever.
		select {}
	},
}

func test2() {
	ipfsURI := "ipfs://c5b0c8a705e0fbd783562ca9680813e064c97a0e8de4f71501601c1d6d1d00b8"
	ipfsHash := extractCIDFromURI(ipfsURI)

	ipfsContent, err := fetchIPFSContent(ipfsHash)
	if err != nil {
		fmt.Println("Error fetching IPFS content:", err)
		return
	}

	fmt.Println("IPFS Content:", string(ipfsContent))
}

func extractCIDFromURI(ipfsURI string) string {
	parts := strings.Split(ipfsURI, "://")
	if len(parts) != 2 || parts[0] != "ipfs" {
		return ""
	}
	return parts[1]
}

func fetchIPFSContent(ipfsHash string) ([]byte, error) {
	url := fmt.Sprintf("https://dweb.link/ipfs/%s", ipfsHash)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch IPFS content. Status code: %d", resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func test1() {
	client, err := ethclient.Dial("https://patient-crimson-slug.matic.discover.quiknode.pro/4cb47dc694ccf2998581548feed08af5477aa84b/")
	if err != nil {
		log.Printf("RPC client err, %s\n", err)
		return
	}
	tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash("0x168b74089afa81a1e19a0d3da37f0d248330ed7cf3c50a4c5e315776b6a9a004"))
	receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash("0x168b74089afa81a1e19a0d3da37f0d248330ed7cf3c50a4c5e315776b6a9a004"))
	logs := receipt.Logs
	status := receipt.Status
	for _, log := range logs {
		fmt.Println(status)
		topics := log.Topics
		fmt.Println(topics[0].String())
		fmt.Println(hex.EncodeToString(log.Data))
	}
	fmt.Println(tx.Data())
}

func testSubscribe() {

	//backend, err := ethclient.Dial("https://blue-cosmopolitan-lambo.discover.quiknode.pro/16de86a8ab8720fb7eb54f3b22e18c78068cf077/")
	//if err != nil {
	//	log.Printf("failed to dial: %v", err)
	//	return
	//}

	rpcCli, err := rpc.Dial("wss://blue-cosmopolitan-lambo.discover.quiknode.pro/16de86a8ab8720fb7eb54f3b22e18c78068cf077/")
	if err != nil {
		log.Printf("failed to dial: %v", err)
		return
	}
	gcli := gethclient.New(rpcCli)

	txch := make(chan common.Hash, 100)
	_, err = gcli.SubscribePendingTransactions(context.Background(), txch)
	if err != nil {
		log.Printf("failed to SubscribePendingTransactions: %v", err)
		return
	}
	for {
		select {
		case txhash := <-txch:
			//tx, _, err := backend.TransactionByHash(context.Background(), txhash)
			//if err != nil {
			//	continue
			//}
			//data, _ := tx.MarshalJSON()
			log.Printf("tx: %v", txhash)
		}
	}
}

func GetAbi() {
	client, err := ethclient.Dial("https://cool-radial-theorem.discover.quiknode.pro/813366a7065d12ddc761101db4934091624c355d/")
	if err != nil {
		return
	}

	var abi1 = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_lido\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_treasury\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requestedBy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ERC20Recovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requestedBy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721Recovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ETHReceived\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LIDO\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TREASURY\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"recoverERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"recoverERC721\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_maxAmount\",\"type\":\"uint256\"}],\"name\":\"withdrawRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"

	abiRes, err := abi.JSON(strings.NewReader(abi1))
	if err != nil {
		log.Fatal(err)
	}
	methods := abiRes.Methods
	for _, method := range methods {
		fmt.Println(method.Sig)
		fmt.Println(method.Name)
		fmt.Println(method.ID)
	}

	url := "https://api.etherscan.io/api?apikey=GG2BR38NQ94SBCTFWTDH72NH53Z9Y1X3IF&module=contract&action=getsourcecode&address=0x388C818CA8B9251b393131C08a736A67ccB19297"

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		fmt.Println("send request fail:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read response:", err)
		return
	}

	fmt.Println(string(body))

	contractAddress := common.HexToAddress("0x388C818CA8B9251b393131C08a736A67ccB19297")

	abi, err := getContractABI(client, contractAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Contract ABI:", abi)
}

func getContractABI(client *ethclient.Client, address common.Address) (interface{}, error) {
	bytecode, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		return "", err
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(bytecode)))
	if err != nil {
		return "", err
	}

	return parsedABI, nil
}
