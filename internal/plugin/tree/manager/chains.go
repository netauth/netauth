package manager

var (
	entityChainConfig = []string{
		"CREATE:plugin-entitycreate",
		"DESTROY:plugin-entitydestroy",
		"SET-SECRET:plugin-presecretchange",
		"SET-SECRET:plugin-postsecretchange",
		"VALIDATE-IDENTITY:plugin-preauthcheck",
		"VALIDATE-IDENTITY:plugin-postauthcheck",
		"LOCK:plugin-entitylock",
		"UNLOCK:plugin-entityunlock",
		"MERGE-METADATA:plugin-entityupdate",
		"UEM-UPSERT:plugin-entityupdate",
		"UEM-CLEARFUZZY:plugin-entityupdate",
		"UEM-CLEAREXACT:plugin-entityupdate",
	}

	groupChainConfig = []string{
		"CREATE:plugin-groupcreate",
		"UPDATE:plugin-groupupdate",
		"UGM-UPSERT:plugin-groupupdate",
		"UGM-CLEARFUZZY:plugin-groupupdate",
		"UGM-CLEAREXACT:plugin-groupupdate",
		"DESTROY:plugin-groupdestroy",
	}
)
