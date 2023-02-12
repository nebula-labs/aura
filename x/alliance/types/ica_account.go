package types

func GetICAAccountOwner(chainId string) (result string) {
	return chainId + "." + "ACCOUNT"
}
