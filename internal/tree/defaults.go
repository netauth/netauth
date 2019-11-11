package tree

var (
	defaultEntityChains = map[string][]string{
		"CREATE": {
			"fail-on-existing-entity",
			"set-entity-id",
			"set-entity-number",
			"set-entity-secret",
			"save-entity",
		},
		"DESTROY": {
			"load-entity",
			"destroy-entity",
		},
		"FETCH": {
			"load-entity",
		},
		"BOOTSTRAP-SERVER": {
			"create-entity-if-missing",
			"ensure-entity-meta",
			"unlock-entity",
			"set-entity-secret",
			"set-entity-capability",
			"save-entity",
		},
		"SET-SECRET": {
			"load-entity",
			"set-entity-secret",
			"save-entity",
		},
		"SET-CAPABILITY": {
			"load-entity",
			"ensure-entity-meta",
			"set-entity-capability",
			"save-entity",
		},
		"DROP-CAPABILITY": {
			"load-entity",
			"ensure-entity-meta",
			"remove-entity-capability",
			"save-entity",
		},
		"ADD-KEY": {
			"load-entity",
			"ensure-entity-meta",
			"add-entity-key",
			"save-entity",
		},
		"DEL-KEY": {
			"load-entity",
			"ensure-entity-meta",
			"del-entity-key",
			"save-entity",
		},
		"VALIDATE-IDENTITY": {
			"load-entity",
			"validate-entity-unlocked",
			"validate-entity-secret",
			"save-entity",
		},
		"MERGE-METADATA": {
			"load-entity",
			"ensure-entity-meta",
			"merge-entity-meta",
			"save-entity",
		},
		"LOCK": {
			"load-entity",
			"ensure-entity-meta",
			"lock-entity",
			"save-entity",
		},
		"UNLOCK": {
			"load-entity",
			"ensure-entity-meta",
			"unlock-entity",
			"save-entity",
		},
		"UEM-UPSERT": {
			"load-entity",
			"ensure-entity-meta",
			"add-untyped-metadata",
			"save-entity",
		},
		"UEM-CLEARFUZZY": {
			"load-entity",
			"ensure-entity-meta",
			"del-untyped-metadata-fuzzy",
			"save-entity",
		},
		"UEM-CLEAREXACT": {
			"load-entity",
			"ensure-entity-meta",
			"del-untyped-metadata-exact",
			"save-entity",
		},
		"GROUP-ADD": {
			"load-entity",
			"ensure-entity-meta",
			"add-direct-group",
			"save-entity",
		},
		"GROUP-DEL": {
			"load-entity",
			"ensure-entity-meta",
			"del-direct-group",
			"save-entity",
		},
	}

	defaultGroupChains = map[string][]string{
		"CREATE": {
			"fail-on-existing-group",
			"set-group-name",
			"set-managing-group",
			"set-group-displayname",
			"set-group-number",
			"save-group",
		},
		"DESTROY": {
			"load-group",
			"destroy-group",
		},
		"FETCH": {
			"load-group",
		},
		"MERGE-METADATA": {
			"load-group",
			"merge-group-meta",
			"save-group",
		},
		"SET-CAPABILITY": {
			"load-group",
			"set-group-capability",
			"save-group",
		},
		"DROP-CAPABILITY": {
			"load-group",
			"remove-group-capability",
			"save-group",
		},
		"UGM-UPSERT": {
			"load-group",
			"add-untyped-metadata",
			"save-group",
		},
		"UGM-CLEARFUZZY": {
			"load-group",
			"del-untyped-metadata-fuzzy",
			"save-group",
		},
		"UGM-CLEAREXACT": {
			"load-group",
			"del-untyped-metadata-exact",
			"save-group",
		},
		"MODIFY-EXPANSIONS": {
			"load-group",
			"check-immediate-expansions",
			"check-expansion-cycles",
			"check-expansion-targets",
			"patch-group-expansions",
			"save-group",
		},
	}
)
