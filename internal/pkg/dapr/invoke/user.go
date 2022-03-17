package invoke

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/Meland-Inc/logservice/internal/pkg/dapr/message"
	daprc "github.com/dapr/go-sdk/client"
)

var userIdAddressCacheMap = map[string]string{
	// "6": strings.ToLower("0x17a243f7Dd13BadE0a7001Ad71a7ef4628A75fCB"),
}

func GetUserWeb3Profile(userId string) (*message.GetUserWeb3ProfileOutput, error) {
	if address, ok := userIdAddressCacheMap[userId]; ok {
		return &message.GetUserWeb3ProfileOutput{
			BlockchainAddress: address,
		}, nil
	}

	client, err := daprc.NewClient()

	if err != nil {
		return nil, err
	}

	input := message.GetUserWeb3ProfileInput{
		UserId: userId,
	}

	inputjson, err := json.Marshal(input)

	if err != nil {
		return nil, err
	}

	content := daprc.DataContent{
		ContentType: "application/json",
		Data:        inputjson,
	}

	output := message.GetUserWeb3ProfileOutput{}
	outputdata, err := client.InvokeMethodWithContent(context.Background(), string(message.AppIdMelandService), string(message.MelandServiceActionGetUserWeb3Profile), "post", &content)

	if err != nil {
		return nil, err
	}

	err = output.UnmarshalJSON(outputdata)

	if err != nil {
		return nil, err
	}

	userIdAddressCacheMap[strings.ToLower(output.BlockchainAddress)] = userId

	return &output, nil
}

var addressUserIdCacheMap = map[string]string{
	// strings.ToLower("0xf1f2eeb0d54ac1f591d294ca95e85478b6324e72"): "3317",
	// strings.ToLower("0x17a243f7Dd13BadE0a7001Ad71a7ef4628A75fCB"): "6",
}

func GetUserIdByBlockchainAddress(blockchainAddress string) (*message.GetUserIdByAddressOutput, error) {
	blockchainAddress = strings.ToLower(blockchainAddress)
	if userId, ok := addressUserIdCacheMap[strings.ToLower(blockchainAddress)]; ok {
		return &message.GetUserIdByAddressOutput{
			UserId: userId,
		}, nil
	}

	client, err := daprc.NewClient()

	if err != nil {
		return nil, err
	}

	input := message.GetUserIdByAddressInput{
		BlockchainAddress: strings.ToLower(blockchainAddress),
	}

	inputjson, err := json.Marshal(input)

	if err != nil {
		return nil, err
	}

	content := daprc.DataContent{
		ContentType: "application/json",
		Data:        inputjson,
	}

	output := message.GetUserIdByAddressOutput{}
	outputdata, err := client.InvokeMethodWithContent(context.Background(), string(message.AppIdMelandService), string(message.MelandServiceActionGetUserIdByAddress), "post", &content)

	if err != nil {
		return nil, err
	}

	addressUserIdCacheMap[output.UserId] = strings.ToLower(input.BlockchainAddress)

	err = output.UnmarshalJSON(outputdata)

	if err != nil {
		return nil, err
	}

	return &output, nil
}
