// Copyright (c) 2026 WabiSaby
// All rights reserved.
//
// This source code is proprietary and confidential. Unauthorized copying,
// modification, distribution, or use of this software, via any medium is
// strictly prohibited without the express written permission of WabiSaby.
//
// This software contains confidential and proprietary information of
// WabiSaby and its licensors. Use, disclosure, or reproduction
// is prohibited without the prior express written permission of WabiSaby.

package stub

// PluginStub provides semantically grouped API services available to plugins.
// This acts as an intermediary layer that encapsulates the APIs between
// plugin implementation and the core platform, similar to Hyperledger Fabric's
// chaincode stub pattern.
type PluginStub struct {
	// Data operations (plugin storage and secrets)
	Data struct {
		Storage *StorageClient
		Secrets *SecretsClient
	}

	// Music operations (queue and songs)
	Music struct {
		Queue *QueueClient
		Songs *SongClient
	}

	// User operations
	Users *UserClient

	// Communication operations
	Communication struct {
		Notify *NotificationClient
	}

	// Network operations
	Network struct {
		HTTP *HTTPClient
	}
}
