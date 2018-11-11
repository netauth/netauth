package tree

var (
	defaultEntityChains = map[string][]string{
		"CREATE": []string{
			"fail-on-existing-entity",
			"set-entity-id",
			"set-entity-number",
			"set-entity-secret",
			"save-entity",
		},
		"DESTROY": []string{
			"destroy-entity",
		},
		"FETCH": []string{
			"load-entity",
		},
		"BOOTSTRAP-SERVER": []string{
			"create-entity-if-missing",
			"ensure-entity-meta",
			"unlock-entity",
			"set-entity-capability",
			"save-entity",
		},
		"SET-SECRET": []string{
			"load-entity",
			"set-entity-secret",
			"save-entity",
		},
		"SET-CAPABILITY": []string{
			"load-entity",
			"ensure-entity-meta",
			"set-entity-capability",
			"save-entity",
		},
		"DROP-CAPABILITY": []string{
			"load-entity",
			"ensure-entity-meta",
			"remove-entity-capability",
			"save-entity",
		},
		"ADD-KEY": []string{
			"load-entity",
			"ensure-entity-meta",
			"add-entity-key",
			"save-entity",
		},
		"DEL-KEY": []string{
			"load-entity",
			"ensure-entity-meta",
			"del-entity-key",
			"save-entity",
		},
		"VALIDATE-IDENTITY": []string{
			"load-entity",
			"validate-entity-unlocked",
			"validate-entity-secret",
			"save-entity",
		},
		"MERGE-METADATA": []string{
			"load-entity",
			"ensure-entity-meta",
			"merge-entity-meta",
			"save-entity",
		},
		"LOCK": []string{
			"load-entity",
			"ensure-entity-meta",
			"lock-entity",
			"save-entity",
		},
		"UNLOCK": []string{
			"load-entity",
			"ensure-entity-meta",
			"unlock-entity",
			"save-entity",
		},
		"UEM-UPSERT": []string{
			"load-entity",
			"ensure-entity-meta",
			"add-untyped-metadata",
			"save-entity",
		},
		"UEM-CLEARFUZZY": []string{
			"load-entity",
			"ensure-entity-meta",
			"del-untyped-metadata-fuzzy",
			"save-entity",
		},
		"UEM-CLEAREXACT": []string{
			"load-entity",
			"ensure-entity-meta",
			"del-untyped-metadata-exact",
			"save-entity",
		},
	}
)
