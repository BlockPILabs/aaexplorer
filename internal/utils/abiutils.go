package utils

var methodMap = map[string]string{
	//---//transfer(address,uint256)： 0xa9059cbb
	"0xa9059cbb": "transfer",
	//---//balanceOf(address)：0x70a08231
	"0x70a08231": "balanceOf",
	//---//decimals()：0x313ce567
	"0x313ce567": "decimals",
	//---//allowance(address,address)： 0xdd62ed3e
	"0xdd62ed3e": "allowance",
	//---//symbol()：0x95d89b41
	"0x95d89b41": "symbol",
	//---//totalSupply()：0x18160ddd
	"0x18160ddd": "totalSupply",
	//---//name()：0x06fdde03
	"0x06fdde03": "name",
	//---//approve(address,uint256)：0x095ea7b3
	"0x095ea7b3": "approve",
	//---//transferFrom(address,address,uint256)： 0x23b872dd
	"0x23b872dd": "transferFrom",
}
