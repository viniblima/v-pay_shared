package enums

const (

	// Wallet
	ActionCreateWallet             = "create_wallet"
	ActionDeleteWalletByID         = "delete_wallet_by_id"
	ActionDeleteAllWalletsByUserID = "delete_all_wallets_by_user_id"

	// Notification
	ActionCreateNotificationTransferSended   = "create_notification_transfer_sended"
	ActionCreateNotificationTransferReceived = "create_notification_transfer_received"

	// Store
	ActionRemoveAllStoresByUserID = "remove_all_stores_by_user_id"
)
